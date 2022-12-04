package main

import (
	"fmt"
	"github.com/kakwa/wows-whaling-simulator/lootbox"
	"os"
)

func main() {
	steel := lootbox.Item{Name: "Steel", Quantity: 1250}
	lootBox, err := lootbox.NewLootBoxFromCSVs([]string{"./rates/droprate-2021-megagift.csv"}, "Santa Mega Gift Containers 2021", 1250, &steel)

	if err != nil {
		fmt.Printf("Error parsing file: %s\n", err.Error())
		os.Exit(1)
	}

	ws, err := lootbox.NewWhalingSession(lootBox, []string{"Warspite", "Poltava"})
	if err != nil {
		fmt.Printf("Error Initializing the whaling Session: %s\n", err.Error())
		os.Exit(1)
	}

	for i := 0; i < 1000; i++ {
		err := ws.Draw()
		if err != nil {
			fmt.Printf("Error Drawing item: %s\n", err.Error())
			os.Exit(1)

		}
	}
	for k, v := range ws.Items {
		fmt.Printf("%6d Item: %s\n", v, k)
	}
}
