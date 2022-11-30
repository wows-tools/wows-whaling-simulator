package main

import (
	"os"
	"fmt"
	"encoding/csv"
	"regexp"
	"strconv"
	mapset "github.com/deckarep/golang-set/v2"
//	"github.com/rdleal/intervalst/interval"
)

type DropRate struct {
	Category string
	Rate  float64
	Pity  uint64
	Quantity uint64
	Items mapset.Set[string]
}

func ParseDropRateFile(filePath string) ([]DropRate, error) {
	item_parsing := regexp.MustCompile(` *| *`)
	var droprates []DropRate
	csvFile, err := os.Open(filePath)
	if err != nil {
    		return nil, err
	}
	defer csvFile.Close()
	 csvLines, err := csv.NewReader(csvFile).ReadAll()
    	if err != nil {
    		return nil, err
    	}
	for _, line := range csvLines {
		category := line[0]
		rate, err := strconv.ParseFloat(line[1], 64)
		if err != nil {
			continue
		}
		quantity, err := strconv.ParseUint(line[2], 10, 64)
		if err != nil {
			continue
		}
		pity, err := strconv.ParseUint(line[3], 10, 64)
		if err != nil {
			continue
		}
		items := item_parsing.Split(line[4], -1)
		dr := DropRate{
			Category: category,
			Rate: rate,
			Quantity: quantity,
			Pity: pity,
			Items: mapset.NewSet[string](),
		}
		for _, item := range items {
			dr.Items.Add(item)
		}
		droprates = append(droprates, dr)
	}
	return droprates, nil
}

func main() {
	_, err := ParseDropRateFile("./rates/droprate-2021-megagift.csv")
	if err != nil {
		fmt.Printf("Error parsing file: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Printf("hello world")
}
