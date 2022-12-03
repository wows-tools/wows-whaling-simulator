package main

import (
	"fmt"
	"github.com/kakwa/wows-whaling-simulator/lootbox"
	"os"
)

func main() {
	steel := lootbox.Item{Name: "steel", Quantity: 1250}
	lootBox, err := lootbox.NewLootBoxFromCSVs([]string{"./rates/droprate-2021-megagift.csv"}, "Santa Mega Gift Containers 2021", 1250, &steel)

	if err != nil {
		fmt.Printf("Error parsing file: %s\n", err.Error())
		os.Exit(1)
	}

	err = lootBox.RefreshDrawTree()

	lootBox.Init()

	if err != nil {
		fmt.Printf("Error building the interval tree: %s\n", err.Error())
		os.Exit(1)
	}

	for i := 0; i < 10000; i++ {
		results, err := lootBox.Draw(false)
		for _, res := range results {
			if err != nil {
				fmt.Printf("Error drawing: %s\n", err.Error())
				os.Exit(1)
			}
			fmt.Printf("%s | %s\n", res.Category.Name, res.Item.Name)
		}
	}
}
