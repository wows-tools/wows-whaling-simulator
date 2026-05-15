package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/kakwa/wows-whaling-simulator/lootbox"
	"github.com/kakwa/wows-whaling-simulator/wows"
	"log"
	"os"
	"strings"
)

func main() {
	num := flag.Int("n", 0, "number of containers opened")
	target := flag.String("target", "", "ship targeted (exclusive with -n")
	realmStr := flag.String("realm", "eu", "Wows realm (eu, na, asia)")
	nick := flag.String("nick", "", "Nickname of the player (requires WOWS_WOWSAPIKEY env var; omit for offline mode)")
	ownedShips := flag.String("owned-ships", "", "Comma-separated list of owned ship names to exclude from drops (offline alternative to -nick)")
	lootboxType := flag.String("lootbox", "", "Lootbox type")
	statsMode := flag.Bool("stats", false, "Enable stats mode (stats on 1000 runs)")
	rawDir := flag.String("raw-dir", "", "Load lootboxes directly from a raw WG vortex JSON directory (bypasses rates/)")
	flag.Parse()

	if len(*target) != 0 && *num != 0 {
		log.Fatal("-n and -target are exclusive flags")
	}

	var lbc map[string]*lootbox.LootBox
	var err error
	if *rawDir != "" {
		lbc, err = lootbox.NewLootBoxCollectionFromVortexDir(*rawDir)
	} else {
		lbc, err = lootbox.NewLootBoxCollection("./rates/")
	}
	if err != nil {
		log.Fatal("Error parsing files: %s\n", err.Error())
	}

	lb, ok := lbc[*lootboxType]
	if !ok {
		fmt.Printf("Unknow lootbox, please chose one with '-lootbox <LB_ID>' in:\n")
		for l := range lbc {
			fmt.Printf("* %s\n", l)
		}
		os.Exit(1)
	}

	var ships []string

	if len(*nick) != 0 {
		api_key := os.Getenv("WOWS_WOWSAPIKEY")
		wowsApi := wows.NewWowsAPI(api_key)
		realm, err := wows.WowsRealm(*realmStr)
		if err != nil {
			log.Fatal("Error Getting the realm %s\n", err.Error())
		}
		err = wowsApi.FillShipMapping()
		if err != nil {
			log.Fatal("Error Prefilling the Ship Mapping: %s\n", err.Error())
		}
		players, err := wowsApi.SearchPlayer(realm, *nick, "exact")
		if err != nil {
			log.Fatal("Error Searching the player: %s\n", err.Error())
		}
		ships, err = wowsApi.GetPlayerShips(realm, *players[0].AccountId)
		if err != nil {
			log.Fatal("Error getting the player's ships: %s\n", err.Error())
		}
	} else if len(*ownedShips) != 0 {
		for _, s := range strings.Split(*ownedShips, ",") {
			s = strings.TrimSpace(s)
			if s != "" {
				ships = append(ships, s)
			}
		}
	}

	wss := lb.NewWhalingStatsSession(ships)

	ws, err := lb.NewWhalingSession(ships)
	if err != nil {
		log.Fatal("Error Initializing the whaling Session: %s\n", err.Error())
	}

	var data []byte
	if *num != 0 && !*statsMode {
		err = ws.SimpleWhaling(*num)
		if err != nil {
			log.Fatal("Error Drawing item: %s\n", err.Error())

		}

		data, err = json.MarshalIndent(ws, "", " ")
		if err != nil {
			log.Fatal(err)
		}
	}
	if len(*target) != 0 && !*statsMode {
		err = ws.TargetWhaling(*target)
		if err != nil {
			log.Fatal("Error Drawing item: %s\n", err.Error())

		}

		data, err = json.MarshalIndent(ws, "", " ")
		if err != nil {
			log.Fatal(err)
		}
	}
	if *statsMode {
		lootbox.InitStatsWorkers()
	}

	if *num != 0 && *statsMode {
		err = wss.StatsSimpleWhaling(*num)
		if err != nil {
			log.Fatal("Error Drawing item: %s\n", err.Error())

		}

		data, err = json.MarshalIndent(wss, "", " ")
		if err != nil {
			log.Fatal(err)
		}

	}

	if len(*target) != 0 && *statsMode {
		err = wss.StatsTargetWhaling(*target)
		if err != nil {
			log.Fatal("Error Drawing item: %s\n", err.Error())

		}

		data, err = json.MarshalIndent(wss, "", " ")
		if err != nil {
			log.Fatal(err)
		}
	}
	if *statsMode {
		lootbox.StopStatsWorkers()
	}
	fmt.Printf(string(data))
}
