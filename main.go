package main

import (
	"fmt"
	"os"
	"github.com/kakwa/wows-whaling-simulator/lootbox"
	"github.com/kakwa/wows-whaling-simulator/wows"
)

func main() {
	fmt.Printf("Launching\n")
	api_key := os.Getenv("WG_API_KEY")
	wowsApi := wows.NewWowsAPI(api_key)
	err := wowsApi.FillShipMapping()
	if err != nil {
		fmt.Printf("Error Prefilling the Ship Mapping: %s\n", err.Error())
		os.Exit(1)
	}
	players, err := wowsApi.SearchPlayer(wows.EURealm, "kakwa", "exact")
	if err != nil {
		fmt.Printf("Error Searching the player: %s\n", err.Error())
		os.Exit(1)
	}

	for _, player := range players {
                fmt.Printf("id %d nick: %s\n", *player.AccountId, *player.Nickname)
        }
	ships, err := wowsApi.GetPlayerShips(wows.EURealm, *players[0].AccountId)
	if err != nil {
		fmt.Printf("Error getting the player's ships: %s\n", err.Error())
		os.Exit(1)
	}
	//for _, ship := range ships {
	//	fmt.Printf("%s\n", ship)
	//}


	lb, err := lootbox.NewLootBoxFromJson("./rates/santa_mega_2021.json")

	if err != nil {
		fmt.Printf("Error parsing file: %s\n", err.Error())
		os.Exit(1)
	}

	ws, err := lootbox.NewWhalingSession(lb, ships)
	if err != nil {
		fmt.Printf("Error Initializing the whaling Session: %s\n", err.Error())
		os.Exit(1)
	}

	for i := 0; i < 120; i++ {
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
