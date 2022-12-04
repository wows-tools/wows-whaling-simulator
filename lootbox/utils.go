package lootbox

import (
	"encoding/csv"
	"os"
	"regexp"
	"strconv"
)

func NewLootBoxFromCSVs(filePaths []string, name string, price float64, compensation *Item) (*LootBox, error) {
	item_parsing := regexp.MustCompile(` *\| *`)
	var lb LootBox
	lb.Name = name
	lb.Price = price

	for _, filePath := range filePaths {
		drop := make(map[string]*ItemCategory)
		csvFile, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		defer csvFile.Close()
		csvLines, err := csv.NewReader(csvFile).ReadAll()
		if err != nil {
			return nil, err
		}

		for i, line := range csvLines {
			if i == 0 {
				continue
			}

			category := line[0]
			rate, err := strconv.ParseFloat(line[1], 64)
			if err != nil {
				return nil, err
			}

			quantity, err := strconv.ParseUint(line[2], 10, 64)
			if err != nil {
				return nil, err
			}

			pityRate, err := strconv.ParseUint(line[3], 10, 64)
			if err != nil {
				return nil, err
			}
			pitiable := false
			if pityRate > 0 {
				lb.Pity = pityRate
				pitiable = true
			}

			transferTo := item_parsing.Split(line[5], -1)
			items := item_parsing.Split(line[6], -1)
			dr := ItemCategory{
				Name:                category,
				DropRate:            rate,
				Pitiable:            pitiable,
				DropRateTransfersTo: transferTo,
				Collectable:         pitiable,
			}

			for _, itemName := range items {
				item := Item{
					Quantity:     quantity,
					Name:         itemName,
					ID:           "undefined",
					Attributes:   make(map[string]string),
					Compensation: compensation,
				}
				if pitiable {
					item.Attributes["tier"] = "V"
				}
				dr.Items = append(dr.Items, &item)
			}
			drop[category] = &dr
		}
		lb.Drops = append(lb.Drops, drop)
	}

	return &lb, nil
}
