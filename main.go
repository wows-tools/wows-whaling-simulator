package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/rdleal/intervalst/interval"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"time"
)

type DropRate struct {
	Category string
	Rate     float64
	Pity     uint64
	Quantity uint64
	Items    mapset.Set[string]
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
		pity, err := strconv.ParseUint(line[3], 10, 64)
		if err != nil {
			return nil, err
		}
		items := item_parsing.Split(line[4], -1)
		dr := DropRate{
			Category: category,
			Rate:     rate,
			Quantity: quantity,
			Pity:     pity,
			Items:    mapset.NewSet[string](),
		}
		for _, item := range items {
			dr.Items.Add(item)
		}
		droprates = append(droprates, dr)
	}

	return droprates, nil
}

func cdrCmp(cumDroprate1, cumDropRate2 float64) int {
	if cumDroprate1 > cumDropRate2 {
		return 1
	}
	if cumDroprate1 == cumDropRate2 {
		return 0
	}
	return -1
}

func BuildSearchTree(dropRates []DropRate) (*interval.SearchTree[DropRate, float64], error) {
	st := interval.NewSearchTree[DropRate](cdrCmp)
	var start float64
	for _, dropRate := range dropRates {
		end := start + dropRate.Rate
		err := st.Insert(start, end, dropRate)
		if err != nil {
			return nil, err
		}
		start = end
	}

	return st, nil
}

func Draw(drTree *interval.SearchTree[DropRate, float64], randSeed *rand.Rand) (*DropRate, error) {
	draw := randSeed.Float64() * 100
	data, ok := drTree.AnyIntersection(draw, draw)
	if !ok {
		return nil, errors.New("did not found a proper interval when drawing")
	}

	return &data, nil
}

func main() {
	dropRates, err := ParseDropRateFile("./rates/droprate-2021-megagift.csv")

	if err != nil {
		fmt.Printf("Error parsing file: %s\n", err.Error())
		os.Exit(1)
	}

	drTree, err := BuildSearchTree(dropRates)

	if err != nil {
		fmt.Printf("Error building the interval tree: %s\n", err.Error())
		os.Exit(1)
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 10000000; i++ {
		data, err := Draw(drTree, r)
		if err != nil {
			fmt.Printf("Error drawing: %s\n", err.Error())
			os.Exit(1)
		}
		fmt.Printf("%s\n", data.Category)
	}
}
