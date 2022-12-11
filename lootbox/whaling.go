package lootbox

import (
	"github.com/barkimedes/go-deepcopy"
	"github.com/kakwa/wows-whaling-simulator/common"
	"math"
)

type WhalingSession struct {
	lootBox          *LootBox
	pityCounter      uint64
	otherItems       map[string]*ItemShortQuantity `json:"other_items"`
	ContainerOpened  uint64                        `json:"container_opened"`
	Pities           uint64                        `json:"pities"`
	Spent            float64                       `json:"game_money_spent"`
	SpentEuro        float64                       `json:"euro_spent"`
	SpentDollar      float64                       `json:"dollar_spent"`
	CollectableItems []*ItemShort                  `json:"collectables_items"`
	OtherItems       []*ItemShortQuantity          `json:"other_items"`
}

func (lb *LootBox) NewWhalingSession(collectables []string) (*WhalingSession, error) {
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
	wlSess.pityCounter = 0
	wlSess.Pities = 0
	wlSess.otherItems = make(map[string]*ItemShortQuantity)
	wlSess.lootBox.Init()
	for _, collectable := range collectables {
		wlSess.lootBox.AddCollectable(collectable)
	}
	wlSess.lootBox.RemoveOwnedCollectablesFromDropableItems()
	wlSess.lootBox.RefreshDrawTree()
	return &wlSess, nil
}

func (ws *WhalingSession) Draw() error {

	usePity := (ws.pityCounter == ws.lootBox.Pity)
	// If we reach pity, reset the counter, otherwise increment it
	if usePity {
		ws.pityCounter = 0
		ws.Pities++
	} else {
		ws.pityCounter++
	}
	results, err := ws.lootBox.Draw(usePity)

	if err != nil {
		return err
	}
	ws.ContainerOpened++
	for _, res := range results {
		if res.Category.Collectable {
			itemShort := ItemShort{
				ID:         res.Item.ID,
				Name:       res.Item.Name,
				Attributes: res.Item.Attributes,
			}
			ws.CollectableItems = append(ws.CollectableItems, &itemShort)
		} else {
			itemShort := ItemShortQuantity{
				ID:         res.Item.ID,
				Name:       res.Item.Name,
				Attributes: res.Item.Attributes,
				Quantity:   res.Item.Quantity,
			}
			if _, ok := ws.otherItems[res.Item.Name]; !ok {
				ws.otherItems[res.Item.Name] = &itemShort
			} else {
				ws.otherItems[res.Item.Name].Quantity += res.Item.Quantity
			}

		}
		// If we drop a pity item, we need to reset the Pity Counter
		if res.Category.Pitiable {
			ws.pityCounter = 0
		}
	}
	return nil
}

func (ws *WhalingSession) Finalize() {
	for _, otherItem := range ws.otherItems {
		ws.OtherItems = append(ws.OtherItems, otherItem)
		ws.Spent = float64(ws.ContainerOpened) * ws.lootBox.Price
		// FIXME Should not hardcode conversion rates here
		ws.SpentEuro = math.Round(ws.Spent/308.641975309*100) / 100
		ws.SpentDollar = math.Round(ws.Spent/280.583613917*100) / 100
	}
}
