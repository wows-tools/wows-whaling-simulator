package lootbox

import (
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/kakwa/wows-whaling-simulator/common"
	"github.com/rdleal/intervalst/interval"
	"math/rand"
	"time"
)

type Item struct {
	Quantity     uint64            `json:"quantity"`     // Quantity dropped
	Name         string            `json:"name"`         // Name of the item
	ID           string            `json:"id"`           // Iternal ID
	Compensation *Item             `json:"compensation"` // Compensation in case the player already have the item
	Attributes   map[string]string `json:"attributes"`   // map of attributes (ex: {"tier": "X"})
}

type ItemShort struct {
	Name       string            `json:"name"`       // Name of the item
	ID         string            `json:"id"`         // Iternal ID
	Attributes map[string]string `json:"attributes"` // map of attributes (ex: {"tier": "X"})
}

type ItemShortQuantity struct {
	Quantity   uint64            `json:"quantity"`   // Quantity of this item dropped
	Name       string            `json:"name"`       // Name of the item
	ID         string            `json:"id"`         // Iternal ID
	Attributes map[string]string `json:"attributes"` // map of attributes (ex: {"tier": "X"})

}

type ItemCategory struct {
	Name                string   `json:"name"`                   // Name of the drop category
	ID                  string   `json:"id"`                     // Iternal ID
	DropRate            float64  `json:"drop_rate"`              // Drop Rate (percentage)
	Collectable         bool     `json:"collectable"`            // If true, If the player already have one of the items, it will be removed from the list of items which can drop
	Pitiable            bool     `json:"pitiable"`               // If true, this item is part of the pity mechanic
	DropRateTransfersTo []string `json:"drop_rate_transfers_to"` // List of Categories the drop rate transfers to if player already has all items
	Items               []*Item  `json:"items"`                  // List of items to get
}

type LootBox struct {
	Name             string                     `json:"name"`
	Img              string                     `json:"img"`
	Weight           int                        `json:"weight"`
	ID               string                     `json:"id"`
	Drops            []map[string]*ItemCategory `json:"drops"`
	Pity             uint64                     `json:"pity"`
	Price            float64                    `json:"price"`
	collectables     mapset.Set[string]
	drawTrees        []*interval.SearchTree[*ItemCategory, float64]
	drawTreesMax     []float64
	drawTreesPity    []*interval.SearchTree[*ItemCategory, float64]
	drawTreesPityMax []float64
	userCollactables mapset.Set[string]
	randSeed         *rand.Rand
}

type DrawResult struct {
	Item     *Item
	Category *ItemCategory
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
	lb.userCollactables.Add(collectable)
}

func (lb *LootBox) IsCollectable(target string) bool {
	return lb.collectables.Contains(target)
}

func (lb *LootBox) Init() {
	lb.userCollactables = mapset.NewSet[string]()
	lb.randSeed = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func (lb *LootBox) RemoveOwnedCollectablesFromDropableItems() {
	for i, drop := range lb.Drops {

		for _, itemCategory := range drop {
			// If the item is not collectable, skip it
			if !itemCategory.Collectable {
				continue
			}
			var cleanedItems []*Item
			var item *Item
			for _, item = range itemCategory.Items {
				// If the item is not in the collectable set, append it
				if !lb.userCollactables.Contains(item.Name) {
					cleanedItems = append(cleanedItems, item)
				}
			}
			// if there is no items, then redistribute the drop rate
			if len(cleanedItems) == 0 {
				lb.Transfer(i, itemCategory, item.Compensation)
			} else {
				itemCategory.Items = cleanedItems
			}
		}
	}
}

func (lb *LootBox) RefreshDrawTree() error {
	// Empty the Draw trees and associated Max boundaries
	lb.drawTrees = nil
	lb.drawTreesMax = nil
	lb.drawTreesPity = nil
	lb.drawTreesPityMax = nil

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

	lb.refreshCollectables()
	return nil
}

func (lb *LootBox) Transfer(index int, cat *ItemCategory, substitute *Item) {
	// remove the entry from the list
	delete(lb.Drops[index], cat.Name)
	spread := 0.0
	// check if the target category is still present
	for _, name := range cat.DropRateTransfersTo {
		if _, ok := lb.Drops[index][name]; ok {
			spread += 1.0
		}
	}

	// If we have nowhere to put our drop rate, then replace the item by the substitute with the original Drop Rate
	if spread == 0 {
		name := "Substitute: " + substitute.Name
		lb.Drops[index][name] = &ItemCategory{
			Name:        name,
			DropRate:    cat.DropRate,
			Collectable: false,
			Pitiable:    true,
			Items:       []*Item{substitute},
		}
	} else {
		// Otherwise, spread the drop
		for _, name := range cat.DropRateTransfersTo {
			if targCat, ok := lb.Drops[index][name]; ok {
				targCat.DropRate += cat.DropRate / spread
			}
		}
	}
}

func (lb *LootBox) Draw(isPity bool) (drawResult []DrawResult, err error) {
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
		// If it's a collectable, get the first item from the "queue"
		if result.Category.Collectable {
			if len(result.Category.Items) > 0 {
				result.Item = result.Category.Items[0]
				// Remove the dropped item from the slice
				result.Category.Items = result.Category.Items[1:]
				lb.AddCollectable(result.Item.Name)

				// If there is no more items, then transfer the drop rate to other Collectables
				if len(result.Category.Items) == 0 {
					lb.Transfer(i, result.Category, result.Item.Compensation)
					// Refresh the Draw tree after that
					lb.RefreshDrawTree()
				}
			} else {
				// This should never happen, the category should be transfered or replaced by the substitute before reaching this state
				return nil, common.ErrNoMoreItem
			}
		} else {
			// If it's not a collectable, get one random item from the list
			draw := lb.randSeed.Intn(len(result.Category.Items))
			result.Item = result.Category.Items[draw]
		}
		drawResult = append(drawResult, result)
	}
	return drawResult, nil
}

func (lb *LootBox) refreshCollectables() {
	collectables := mapset.NewSet[string]()
	for _, drop := range lb.Drops {
		for _, itemCategory := range drop {
			// If the item is not collectable, skip it
			if !itemCategory.Collectable {
				continue
			}
			for _, item := range itemCategory.Items {
				collectables.Add(item.Name)
			}
		}
	}
	lb.collectables = collectables
}

func (lb *LootBox) ListCollectables() (ret []*ItemShort) {
	for _, drop := range lb.Drops {
		for _, itemCategory := range drop {
			// If the item is not collectable, skip it
			if !itemCategory.Collectable {
				continue
			}
			for _, item := range itemCategory.Items {
				itemShort := ItemShort{
					ID:         item.ID,
					Name:       item.Name,
					Attributes: item.Attributes,
				}

				ret = append(ret, &itemShort)
			}
		}
	}
	return ret
}
