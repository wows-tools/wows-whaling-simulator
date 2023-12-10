package main

import (
	"encoding/json"
	"fmt"
	"github.com/kakwa/wows-whaling-simulator/lootbox"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func main() {
	//urls := CollectLootboxURLs()
	//for i, url := range urls {
	//    fmt.Printf("curl '%s' | jq . >crap/%d.json\n", url ,i)
	//}

	// Open the JSON file
	file, err := ioutil.ReadFile("crap/1.json")
	if err != nil {
		fmt.Println("Error reading the file:", err)
		return
	}

	// Define a struct to hold the JSON data
	var lb WgLootbox

	// Unmarshal the JSON data into the struct
	err = json.Unmarshal(file, &lb)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	lbws, err := convert(&lb)
	if err != nil {
		fmt.Println("Error converting JSON:", err)
		return
	}

	jsonData, err := json.MarshalIndent(lbws, "", "    ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	// Print the pretty JSON
	fmt.Println(string(jsonData))
}

func name2id(input string) string {
	// Convert to lowercase
	result := strings.ToLower(input)

	// Replace spaces with underscores
	result = strings.ReplaceAll(result, " ", "_")

	return result
}

func getWeight(name string) int {
	// TODO switch case on name pattern (ex: contains "premium")
	return 40
}

func convert(lb *WgLootbox) (*lootbox.LootBox, error) {
	var lbws lootbox.LootBox
	lbws.Name = lb.Data.Title
	lbws.ID = name2id(lb.Data.Title)
	lbws.Weight = getWeight(lbws.ID)
	lbws.Img = "/resources/" + filepath.Base(lb.Data.Icons.Large)

	lbws.ExchangeRateEuro = 303.96
	lbws.ExchangeRateDollar = 289.31

	return &lbws, nil
}
