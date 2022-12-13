package main

import (
	"flag"
	"fmt"
	"github.com/kakwa/wows-whaling-simulator/lootbox"
	"github.com/kakwa/wows-whaling-simulator/wows"
	"os"
)

func main() {
	num := flag.Int("n", 10, "number of containers opened")
	realmStr := flag.String("realm", "eu", "Wows realm (eu, na, asia)")
	nick := flag.String("nick", "", "Nickname of the player")
	lootboxType := flag.String("lootbox", "", "Lootbox type")
	flag.Parse()

	lbc, err := lootbox.NewLootBoxCollection("./rates/")
	if err != nil {
		fmt.Printf("Error parsing files: %s\n", err.Error())
		os.Exit(1)
	}

	lb, ok := lbc[*lootboxType]
	if !ok {
		fmt.Printf("Unknow lootbox, please chose one with '-lootbox <LB_ID>' in:\n")
		for l := range lbc {
			fmt.Printf("* %s\n", l)
		}
		os.Exit(1)
	}

	api_key := os.Getenv("WOWS_WOWSAPIKEY")
	wowsApi := wows.NewWowsAPI(api_key)
	realm, err := wows.WowsRealm(*realmStr)
	if err != nil {
		fmt.Printf("Error Getting the realm %s\n", err.Error())
		os.Exit(1)
	}
	err = wowsApi.FillShipMapping()
	if err != nil {
		fmt.Printf("Error Prefilling the Ship Mapping: %s\n", err.Error())
		os.Exit(1)
	}
	players, err := wowsApi.SearchPlayer(realm, *nick, "exact")
	if err != nil {
		fmt.Printf("Error Searching the player: %s\n", err.Error())
		os.Exit(1)
	}

	ships, err := wowsApi.GetPlayerShips(realm, *players[0].AccountId)
	if err != nil {
		fmt.Printf("Error getting the player's ships: %s\n", err.Error())
		os.Exit(1)
	}

	ws, err := lb.NewWhalingSession(ships)
	if err != nil {
		fmt.Printf("Error Initializing the whaling Session: %s\n", err.Error())
		os.Exit(1)
	}

	for i := 0; i < *num; i++ {
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
