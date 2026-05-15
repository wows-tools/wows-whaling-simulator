package main

import (
	"path/filepath"
	"strings"

	"github.com/kakwa/wows-whaling-simulator/lootbox"
)

// ShipInfo is kept for the optional WG API ship-name enrichment path (-wg-key).
// Currently unused as ship names and tiers come directly from additionalData in
// the raw vortex JSON.
type ShipInfo struct {
	Name      string
	Tier      int
	IsPremium bool
	IsSpecial bool
}

// convert parses raw vortex JSON bytes into a *lootbox.LootBox.
func convert(raw []byte, _ map[int]ShipInfo) (*lootbox.LootBox, error) {
	return lootbox.NewLootBoxFromVortexBytes(raw)
}

// iconImg returns the /resources/<basename> path for the best available icon.
func iconImg(icons Icon) string {
	for _, u := range []string{icons.Large, icons.Default, icons.Small} {
		if u == "" {
			continue
		}
		base := filepath.Base(u)
		if idx := strings.Index(base, "?"); idx != -1 {
			base = base[:idx]
		}
		return "/resources/" + base
	}
	return ""
}
