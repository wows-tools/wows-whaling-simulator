package main

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

type ShipInfo struct {
	Name      string
	Tier      int
	IsPremium bool
	IsSpecial bool
}

// Output structs mirror lootbox.LootBox JSON shape so the rates files are drop-in compatible.

type OutLootBox struct {
	Name               string                         `json:"name"`
	Img                string                         `json:"img"`
	ID                 string                         `json:"id"`
	Weight             int                            `json:"weight"`
	ExchangeRateEuro   float64                        `json:"exchange_rate_euro"`
	ExchangeRateDollar float64                        `json:"exchange_rate_dollar"`
	Drops              []map[string]*OutItemCategory  `json:"drops"`
	Pity               uint64                         `json:"pity"`
	Price              float64                        `json:"price"`
}

type OutItemCategory struct {
	Name                string     `json:"name"`
	ID                  string     `json:"id"`
	DropRate            float64    `json:"drop_rate"`
	Collectable         bool       `json:"collectable"`
	Pitiable            bool       `json:"pitiable"`
	DropRateTransfersTo []string   `json:"drop_rate_transfers_to"`
	Items               []*OutItem `json:"items"`
}

type OutItem struct {
	Quantity     uint64            `json:"quantity"`
	Name         string            `json:"name"`
	ID           string            `json:"id"`
	Compensation *OutItem          `json:"compensation,omitempty"`
	Attributes   map[string]string `json:"attributes"`
}

// rewardTypeMeta maps WG reward types to human-readable names and attribute hints.
var rewardTypeMeta = map[string]struct {
	name       string
	id         string
	attrType   string
	resourceID string
}{
	"signal":          {"Signals", "signals", "signal", ""},
	"credits":         {"Credits", "credits", "resource", "credits"},
	"doubloons":       {"Doubloons", "doubloons", "resource", "doubloons"},
	"coal":            {"Coal", "coal", "resource", "coal"},
	"steel":           {"Steel", "steel", "resource", "steel"},
	"elite_xp":        {"Elite Commander XP", "elite_commander_xp", "resource", "elite_commander_xp"},
	"premium_account": {"Premium Days", "premium_days", "resource", "premium_days"},
	"port_slot":       {"Port Slot", "port_slot", "resource", "port_slot"},
	"bonus":           {"Economic Bonuses", "economic_bonus", "economic bonus", ""},
	"camouflage":      {"Camouflages", "camouflage", "camo", ""},
	"xp":              {"XP", "xp", "resource", "xp"},
	"free_xp":         {"Free XP", "free_xp", "resource", "free_xp"},
	"captain_xp":      {"Commander XP", "commander_xp", "resource", "commander_xp"},
}

func tierToRoman(tier int) string {
	switch tier {
	case 1:
		return "I"
	case 2:
		return "II"
	case 3:
		return "III"
	case 4:
		return "IV"
	case 5:
		return "V"
	case 6:
		return "VI"
	case 7:
		return "VII"
	case 8:
		return "VIII"
	case 9:
		return "IX"
	case 10:
		return "X"
	default:
		return strconv.Itoa(tier)
	}
}

func iconImg(icons Icon) string {
	for _, u := range []string{icons.Large, icons.Default, icons.Small} {
		if u == "" {
			continue
		}
		base := filepath.Base(u)
		if idx := strings.Index(base, "?"); idx != -1 {
			base = base[:idx]
		}
		return "/resources/" + base
	}
	return ""
}

func fillerItem(f Filler) *OutItem {
	typeName := f.Type
	id := ""
	name := typeName
	if meta, ok := rewardTypeMeta[typeName]; ok {
		name = meta.name
		id = meta.resourceID
	}
	attrs := map[string]string{"type": "resource"}
	if id != "" {
		attrs["resource_id"] = id
	}
	return &OutItem{
		Quantity:   uint64(f.Amount),
		Name:       name,
		ID:         id,
		Attributes: attrs,
	}
}

func convert(lb *WgLootbox, ships map[int]ShipInfo) (*OutLootBox, error) {
	id := name2id(lb.Data.Title)

	var pity uint64
	if lb.Data.SavePoint != nil {
		pity = uint64(*lb.Data.SavePoint)
	}

	compensation := fillerItem(lb.Data.Filler)

	out := &OutLootBox{
		Name:               lb.Data.Title,
		Img:                iconImg(lb.Data.Icons),
		ID:                 id,
		Weight:             getWeight(id),
		ExchangeRateEuro:   303.96,
		ExchangeRateDollar: 289.31,
		Pity:               pity,
		Price:              250,
	}

	for _, slot := range lb.Data.Slots {
		dropMap := make(map[string]*OutItemCategory)

		// Common (non-ship) rewards — group by type.
		for _, reward := range slot.CommonRewards.List0.Rewards {
			meta, known := rewardTypeMeta[reward.Type]
			catName := reward.Type
			catID := reward.Type
			attrType := reward.Type
			resourceID := ""
			if known {
				catName = meta.name
				catID = meta.id
				attrType = meta.attrType
				resourceID = meta.resourceID
			}

			attrs := map[string]string{"type": attrType}
			if resourceID != "" {
				attrs["resource_id"] = resourceID
			}

			// For signals/bonus items, try to get a more specific name from the ID field.
			itemName := catName
			itemID := name2id(catID)
			switch v := reward.ID.(type) {
			case float64:
				itemID = strconv.Itoa(int(v))
			case string:
				if v != "" {
					itemName = v
					itemID = name2id(v)
				}
			}

			item := &OutItem{
				Quantity:   uint64(reward.Amount),
				Name:       itemName,
				ID:         itemID,
				Attributes: attrs,
			}

			if cat, exists := dropMap[catName]; exists {
				cat.Items = append(cat.Items, item)
			} else {
				dropMap[catName] = &OutItemCategory{
					Name:                catName,
					ID:                  catID,
					DropRate:            reward.ProbabilityDisplayed,
					Collectable:         false,
					Pitiable:            false,
					DropRateTransfersTo: []string{},
					Items:               []*OutItem{item},
				}
			}
		}

		// Valuable (ship) rewards — one category per ValuableRewardList.
		var valCatNames []string
		valCats := make(map[string]*OutItemCategory)

		for _, valList := range slot.ValuableRewards {
			catName := valList.Title
			if catName == "" {
				catName = valList.ShortTitle
			}
			if catName == "" {
				catName = "Ships"
			}

			var pityVal uint64
			if valList.SavePoint > 0 {
				pityVal = uint64(valList.SavePoint)
				if out.Pity == 0 || pityVal < out.Pity {
					out.Pity = pityVal
				}
			}

			cat := &OutItemCategory{
				Name:                catName,
				ID:                  name2id(catName),
				DropRate:            valList.ProbabilityDisplayed,
				Collectable:         valList.HasUniqueRewards,
				Pitiable:            true,
				DropRateTransfersTo: []string{},
				Items:               []*OutItem{},
			}

			for _, r := range valList.Rewards {
				if r.Type != "ship" {
					continue
				}

				shipName := fmt.Sprintf("ship_%d", r.ID)
				tierStr := ""
				isRare := false

				if info, ok := ships[r.ID]; ok {
					shipName = info.Name
					tierStr = tierToRoman(info.Tier)
					isRare = info.IsSpecial || info.IsPremium
				}

				attrs := map[string]string{
					"type": "ship",
					"rare": strconv.FormatBool(isRare),
				}
				if tierStr != "" {
					attrs["tier"] = tierStr
				}

				cat.Items = append(cat.Items, &OutItem{
					Quantity:     1,
					Name:         shipName,
					ID:           "undefined",
					Compensation: compensation,
					Attributes:   attrs,
				})
			}

			if len(cat.Items) > 0 {
				valCatNames = append(valCatNames, catName)
				valCats[catName] = cat
			}
		}

		// Each ship category transfers its drop rate to all other ship categories.
		for i, name := range valCatNames {
			var others []string
			for j, other := range valCatNames {
				if i != j {
					others = append(others, other)
				}
			}
			valCats[name].DropRateTransfersTo = others
			dropMap[name] = valCats[name]
		}

		out.Drops = append(out.Drops, dropMap)
	}

	return out, nil
}
