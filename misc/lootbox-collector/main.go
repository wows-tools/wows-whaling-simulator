package main

import (
	"fmt"
)

func main() {
    urls := CollectLootboxURLs()
    for i, url := range urls {
        fmt.Printf("%d %s\n",i, url)
    }
}
