package lootbox

import (
	"github.com/barkimedes/go-deepcopy"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/kakwa/wows-whaling-simulator/common"
	"github.com/rdleal/intervalst/interval"
	"math/rand"
	"time"
)

type Item struct {
	Quantity     uint64 `json:"quantity"`     // Quantity dropped
	Name         string `json:"name"`         // Name of the item
	ID           string `json:"id"`           // Iternal ID
	Compensation *Item  `json:"compensation"` // Compensation in case the player already have the item
}

type ItemCategory struct {
	Name                string   `json:"name"`                   // Name of the drop category
	DropRate            float64  `json:"drop_rate"`              // Drop Rate (percentage)
	Collectable         bool     `json:"collectable"`            // If true, If the player already have one of the items, it will be removed from the list of items which can drop
	Pitiable            bool     `json:"pitiable"`               // If true, this item is part of the pity mechanic
	DropRateTransfersTo []string `json:"drop_rate_transfers_to"` // List of Categories the drop rate transfers to if player already has all items
	Items               []Item   `json:"items"`                  // List of items to get
}

type LootBox struct {
	Name             string            `json:"name"`
	Drops            [][]*ItemCategory `json:"drops"`
	Pity             uint64            `json:"pity"`
	Price            float64           `json:"price"`
	drawTrees        []*interval.SearchTree[*ItemCategory, float64]
	drawTreesMax     []float64
	drawTreesPity    []*interval.SearchTree[*ItemCategory, float64]
	drawTreesPityMax []float64
	collectable      mapset.Set[string]
	randSeed         *rand.Rand
}

type DrawResult struct {
	Item     *Item
	Category *ItemCategory
}

type WhalingSession struct {
	lootBox          *LootBox
	CollectableItems []Item
	ContainerOpened  uint64
	PityCounter      uint64
	Items            map[string]uint64
}

func NewWhalingSession(lb *LootBox, collectable []string) (*WhalingSession, error) {
	var ok bool
	dataLb, err := deepcopy.Anything(lb)
	if err != nil {
		return nil, err
	}
	var wlSess WhalingSession
	wlSess.lootBox, ok = dataLb.(*LootBox)
	if !ok {
		return nil, common.ErrCastCopyLB
	}
	wlSess.ContainerOpened = 0
	wlSess.PityCounter = 0
	wlSess.Items = make(map[string]uint64)
	wlSess.lootBox.Init()
	wlSess.lootBox.RefreshDrawTree()
	return &wlSess, nil
}

func (ws *WhalingSession) Draw() error {

	usePity := (ws.PityCounter == ws.lootBox.Pity)
	// If we reach pity, reset the counter, otherwise increment it
	if usePity {
		ws.PityCounter = 0
	} else {
		ws.PityCounter++
	}
	results, err := ws.lootBox.Draw(usePity)

	if err != nil {
		return err
	}
	ws.ContainerOpened++
	for _, res := range results {
		if _, ok := ws.Items[res.Item.Name]; !ok {
			ws.Items[res.Item.Name] = 0
		}
		ws.Items[res.Item.Name] += res.Item.Quantity
		// If we drop a pity item, we need to reset the Pity Counter
		if res.Category.Pitiable {
			ws.PityCounter = 0
		}
	}
	return nil
}

func cdrCmp(cumDroprate1, cumItemCategory2 float64) int {
	if cumDroprate1 > cumItemCategory2 {
		return 1
	}
	if cumDroprate1 == cumItemCategory2 {
		return 0
	}
	return -1
}

func (lb *LootBox) AddCollectable(collectable string) {
	lb.collectable.Add(collectable)
}

func (lb *LootBox) Init() {
	lb.collectable = mapset.NewSet[string]()
	lb.randSeed = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func (lb *LootBox) RefreshDrawTree() error {
	// TODO handle Collectables
	for _, drop := range lb.Drops {
		var start float64
		var end float64
		var startPity float64
		var endPity float64
		st := interval.NewSearchTree[*ItemCategory](cdrCmp)
		stPity := interval.NewSearchTree[*ItemCategory](cdrCmp)
		for _, itemCategory := range drop {
			// Shuffle the items around, to not get the same sequence of drops for collectable
			rand.Shuffle(len(itemCategory.Items), func(i, j int) {
				itemCategory.Items[i], itemCategory.Items[j] = itemCategory.Items[j], itemCategory.Items[i]
			})

			// Insert the category in the main draw tree
			end = start + itemCategory.DropRate
			err := st.Insert(start, end, itemCategory)
			if err != nil {
				return err
			}
			start = end

			// If the item category is not pitiable, don't add it to the pity interval tree
			if !itemCategory.Pitiable {
				continue
			}
			// Insert the item category in the pity draw tree
			endPity = startPity + itemCategory.DropRate
			err = stPity.Insert(startPity, endPity, itemCategory)
			if err != nil {
				return err
			}
			startPity = endPity
		}
		lb.drawTrees = append(lb.drawTrees, st)
		lb.drawTreesMax = append(lb.drawTreesMax, end)
		lb.drawTreesPity = append(lb.drawTreesPity, stPity)
		lb.drawTreesPityMax = append(lb.drawTreesPityMax, endPity)
	}

	return nil
}

func (lb *LootBox) Draw(isPity bool) (drawResult []DrawResult, err error) {
	FakeItem := Item{Name: "Fake Item", Quantity: 1}
	var ok bool
	for i, _ := range lb.drawTrees {
		var result DrawResult
		if !isPity {
			draw := lb.randSeed.Float64() * lb.drawTreesMax[i]
			result.Category, ok = lb.drawTrees[i].AnyIntersection(draw, draw)
		} else {
			draw := lb.randSeed.Float64() * lb.drawTreesPityMax[i]
			result.Category, ok = lb.drawTreesPity[i].AnyIntersection(draw, draw)
		}
		if !ok {
			return nil, common.ErrDraw
		}
		if result.Category.Collectable {
			// TODO FIXME, recompute and dispatch drop rate in case the last item of this category has dropped
			if len(result.Category.Items) > 0 {
				result.Item = &result.Category.Items[0]
				// Remove the dropped item from the slice
				result.Category.Items = result.Category.Items[1:]
			} else {
				result.Item = &FakeItem
			}
			lb.AddCollectable(result.Item.Name)
		} else {
			draw := lb.randSeed.Intn(len(result.Category.Items))
			result.Item = &result.Category.Items[draw]
		}
		drawResult = append(drawResult, result)
	}
	return drawResult, nil
}
