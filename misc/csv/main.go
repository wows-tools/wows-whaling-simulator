package main

import (
	"encoding/json"
	"github.com/kakwa/wows-whaling-simulator/lootbox"
	"io/ioutil"
	"log"
)

func main() {
	steel := lootbox.Item{Name: "steel", Quantity: 1250}
	lootBox, err := lootbox.NewLootBoxFromCSVs([]string{"./rates/droprate-2021-megagift.csv"}, "Santa Mega Gift Containers 2021", 1250, &steel)
	if err != nil {
		log.Fatal(err)
	}
	file, err := json.MarshalIndent(lootBox, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile("test.json", file, 0644)
	if err != nil {
		log.Fatal(err)
	}

}
