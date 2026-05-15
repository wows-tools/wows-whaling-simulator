package lootbox

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// ---------------------------------------------------------------------------
// Raw WG vortex API structs (minimal — only fields we need)
// ---------------------------------------------------------------------------

type vortexEnvelope struct {
	Status string      `json:"status"`
	Data   vortexData  `json:"data"`
}

type vortexData struct {
	ID        int            `json:"id"`
	Name      string         `json:"name"`
	Title     string         `json:"title"`
	IsPremium bool           `json:"isPremium"`
	SavePoint *int           `json:"savePoint"`
	Filler    vortexFiller   `json:"filler"`
	Icons     vortexIcon     `json:"icons"`
	Slots     []vortexSlot   `json:"slots"`
}

type vortexFiller struct {
	Amount int         `json:"amount"`
	Type   string      `json:"type"`
	ID     interface{} `json:"id"`
}

type vortexIcon struct {
	Small   string `json:"small"`
	Large   string `json:"large"`
	Default string `json:"default"`
}

type vortexSlot struct {
	CommonRewards   vortexCommonRewards            `json:"commonRewards"`
	ValuableRewards map[string]vortexValuableList  `json:"valuableRewards"`
}

type vortexCommonRewards struct {
	List0 vortexRewardPool `json:"list0"`
}

type vortexRewardPool struct {
	Rewards []vortexCommonReward `json:"rewards"`
}

type vortexCommonReward struct {
	Amount               int         `json:"amount"`
	Type                 string      `json:"type"`
	ID                   interface{} `json:"id"`
	Weight               int         `json:"weight"`
	ProbabilityDisplayed float64     `json:"probabilityDisplayed"`
}

type vortexValuableList struct {
	SavePoint            int                    `json:"savePoint"`
	Rewards              []vortexValuableReward `json:"rewards"`
	ProbabilityDisplayed float64                `json:"probabilityDisplayed"`
	HasUniqueRewards     bool                   `json:"hasUniqueRewards"`
	Title                string                 `json:"title"`
	ShortTitle           string                 `json:"shortTitle"`
}

type vortexValuableReward struct {
	Amount         int                    `json:"amount"`
	Type           string                 `json:"type"`
	ID             int                    `json:"id"`
	AdditionalData *vortexAdditionalData  `json:"additionalData"`
}

type vortexAdditionalData struct {
	Title     string `json:"title"`
	Level     int    `json:"level"`
	IsPremium bool   `json:"isPremium"`
	IsSpecial bool   `json:"isSpecial"`
}

// ---------------------------------------------------------------------------
// Metadata for common reward types
// ---------------------------------------------------------------------------

type rewardMeta struct {
	name       string
	id         string
	attrType   string
	resourceID string
}

var vortexTypeMeta = map[string]rewardMeta{
	"signal":             {"Signals", "signals", "signal", ""},
	"credits":            {"Credits", "credits", "resource", "credits"},
	"gold":               {"Doubloons", "doubloons", "resource", "doubloons"},
	"coal":               {"Coal", "coal", "resource", "coal"},
	"steel":              {"Steel", "steel", "resource", "steel"},
	"elite_xp":           {"Elite Commander XP", "elite_commander_xp", "resource", "elite_commander_xp"},
	"wows_premium":       {"Premium Days", "premium_days", "resource", "premium_days"},
	"slots":              {"Port Slot", "port_slot", "resource", "port_slot"},
	"camoboost":          {"Economic Bonuses", "economic_bonus", "economic bonus", ""},
	"camouflage":         {"Camouflages", "camouflage", "camo", ""},
	"permoflage":         {"Permanent Camouflages", "permoflage", "camo", ""},
	"style":              {"Styles", "style", "camo", ""},
	"free_xp":            {"Free XP", "free_xp", "resource", "free_xp"},
	"xp":                 {"XP", "xp", "resource", "xp"},
	"collection_album":   {"Collection Album", "collection_album", "resource", "collection_album"},
	"lootbox":            {"Container", "container", "resource", ""},
	"shipdestruction":    {"Supercontainer", "supercontainer", "resource", ""},
	"recruitment_points": {"Recruitment Points", "recruitment_points", "resource", "recruitment_points"},
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func vortexTierRoman(tier int) string {
	switch tier {
	case 1:  return "I"
	case 2:  return "II"
	case 3:  return "III"
	case 4:  return "IV"
	case 5:  return "V"
	case 6:  return "VI"
	case 7:  return "VII"
	case 8:  return "VIII"
	case 9:  return "IX"
	case 10: return "X"
	default: return strconv.Itoa(tier)
	}
}

func vortexName2ID(s string) string {
	s = strings.ToLower(s)
	for _, q := range []string{"'", "’", "“", "”", "«", "»", "\""} {
		s = strings.ReplaceAll(s, q, "")
	}
	var b strings.Builder
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == ' ' || r == '_' || r == '-':
			b.WriteByte('_')
		}
	}
	return b.String()
}

func vortexGetWeight(id string) int {
	switch {
	case strings.Contains(id, "ultra"):
		return 20
	case strings.Contains(id, "mega"):
		return 40
	case strings.Contains(id, "big"):
		return 110
	default:
		return 1100
	}
}

func vortexIconImg(icons vortexIcon) string {
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

func vortexFillerItem(f vortexFiller) *Item {
	name := f.Type
	resID := ""
	if meta, ok := vortexTypeMeta[f.Type]; ok {
		name = meta.name
		resID = meta.resourceID
	}
	attrs := map[string]string{"type": "resource"}
	if resID != "" {
		attrs["resource_id"] = resID
	}
	return &Item{
		Quantity:   uint64(f.Amount),
		Name:       name,
		ID:         resID,
		Attributes: attrs,
	}
}

// ---------------------------------------------------------------------------
// Conversion
// ---------------------------------------------------------------------------

// groupKey identifies a pool of common rewards that share type and probability.
type groupKey struct {
	typ  string
	prob float64
}

func vortexConvert(raw *vortexEnvelope) (*LootBox, error) {
	d := &raw.Data
	id := vortexName2ID(d.Title)

	var pity uint64
	if d.SavePoint != nil {
		pity = uint64(*d.SavePoint)
	}

	compensation := vortexFillerItem(d.Filler)

	lb := &LootBox{
		Name:               d.Title,
		Img:                vortexIconImg(d.Icons),
		ID:                 id,
		Weight:             vortexGetWeight(id),
		ExchangeRateEuro:   303.96,
		ExchangeRateDollar: 289.31,
		Pity:               pity,
		Price:              250,
	}

	for _, slot := range d.Slots {
		dropMap := make(map[string]*ItemCategory)

		// --- Common (non-ship) rewards ---
		// Group by (type, probabilityDisplayed) so items with equal weight
		// form one category and the category drop rate = count × item_prob.
		type groupEntry struct {
			key   groupKey
			items []*Item
		}
		groupOrder := []groupKey{}
		groups := map[groupKey]*groupEntry{}

		for _, r := range slot.CommonRewards.List0.Rewards {
			k := groupKey{r.Type, r.ProbabilityDisplayed}
			if _, exists := groups[k]; !exists {
				groupOrder = append(groupOrder, k)
				groups[k] = &groupEntry{key: k}
			}

			meta, known := vortexTypeMeta[r.Type]
			attrType := r.Type
			resID := ""
			if known {
				attrType = meta.attrType
				resID = meta.resourceID
			}

			attrs := map[string]string{"type": attrType}
			if resID != "" {
				attrs["resource_id"] = resID
			}

			itemID := resID
			itemName := ""
			if known {
				itemName = meta.name
			} else {
				itemName = r.Type
			}
			switch v := r.ID.(type) {
			case float64:
				itemID = strconv.Itoa(int(v))
			case string:
				if v != "" {
					itemID = vortexName2ID(v)
					itemName = v
				}
			}

			groups[k].items = append(groups[k].items, &Item{
				Quantity:   uint64(r.Amount),
				Name:       itemName,
				ID:         itemID,
				Attributes: attrs,
			})
		}

		// Count how many distinct probability groups exist per type (for naming).
		typeGroupCount := map[string]int{}
		for _, k := range groupOrder {
			typeGroupCount[k.typ]++
		}
		typeGroupIdx := map[string]int{}

		for _, k := range groupOrder {
			g := groups[k]
			meta, known := vortexTypeMeta[k.typ]
			baseName := k.typ
			if known {
				baseName = meta.name
			}
			catID := vortexName2ID(baseName)

			// Disambiguate when the same type appears at multiple probabilities.
			name := baseName
			if typeGroupCount[k.typ] > 1 {
				typeGroupIdx[k.typ]++
				name = fmt.Sprintf("%s %d", baseName, typeGroupIdx[k.typ])
				catID = fmt.Sprintf("%s_%d", catID, typeGroupIdx[k.typ])
			}

			dropRate := float64(len(g.items)) * k.prob

			if cat, exists := dropMap[name]; exists {
				cat.Items = append(cat.Items, g.items...)
				cat.DropRate += dropRate
			} else {
				dropMap[name] = &ItemCategory{
					Name:                name,
					ID:                  catID,
					DropRate:            dropRate,
					Collectable:         false,
					Pitiable:            false,
					DropRateTransfersTo: []string{},
					Items:               g.items,
				}
			}
		}

		// --- Valuable (ship) rewards ---
		var shipCatNames []string
		shipCats := map[string]*ItemCategory{}

		for _, vlist := range slot.ValuableRewards {
			catName := vlist.Title
			if catName == "" {
				catName = vlist.ShortTitle
			}
			if catName == "" {
				catName = "Ships"
			}

			// Use list's savePoint to refine pity if tighter.
			if vlist.SavePoint > 0 {
				sp := uint64(vlist.SavePoint)
				if lb.Pity == 0 || sp < lb.Pity {
					lb.Pity = sp
				}
			}

			cat := &ItemCategory{
				Name:                catName,
				ID:                  vortexName2ID(catName),
				DropRate:            vlist.ProbabilityDisplayed,
				Collectable:         vlist.HasUniqueRewards,
				Pitiable:            true,
				DropRateTransfersTo: []string{},
				Items:               []*Item{},
			}

			for _, r := range vlist.Rewards {
				if r.Type != "ship" {
					continue
				}

				shipName := fmt.Sprintf("ship_%d", r.ID)
				tierStr := ""
				isRare := false

				if r.AdditionalData != nil {
					if r.AdditionalData.Title != "" {
						shipName = r.AdditionalData.Title
					}
					if r.AdditionalData.Level > 0 {
						tierStr = vortexTierRoman(r.AdditionalData.Level)
					}
					isRare = r.AdditionalData.IsPremium || r.AdditionalData.IsSpecial
				}

				attrs := map[string]string{
					"type": "ship",
					"rare": strconv.FormatBool(isRare),
				}
				if tierStr != "" {
					attrs["tier"] = tierStr
				}

				cat.Items = append(cat.Items, &Item{
					Quantity:     1,
					Name:         shipName,
					ID:           "undefined",
					Compensation: compensation,
					Attributes:   attrs,
				})
			}

			if len(cat.Items) > 0 {
				shipCatNames = append(shipCatNames, catName)
				shipCats[catName] = cat
			}
		}

		// Each ship category transfers to all other ship categories.
		for i, name := range shipCatNames {
			var others []string
			for j, other := range shipCatNames {
				if i != j {
					others = append(others, other)
				}
			}
			shipCats[name].DropRateTransfersTo = others
			dropMap[name] = shipCats[name]
		}

		lb.Drops = append(lb.Drops, dropMap)
	}

	return lb, nil
}

// ---------------------------------------------------------------------------
// Public API
// ---------------------------------------------------------------------------

// NewLootBoxFromVortexBytes parses a raw WG vortex API JSON response into a LootBox.
func NewLootBoxFromVortexBytes(data []byte) (*LootBox, error) {
	var env vortexEnvelope
	if err := json.Unmarshal(data, &env); err != nil {
		return nil, err
	}
	if env.Status != "ok" {
		return nil, fmt.Errorf("vortex response status %q", env.Status)
	}
	return vortexConvert(&env)
}

// NewLootBoxCollectionFromVortexDir reads a directory of raw vortex JSON files
// and returns a map keyed by lootbox ID.
func NewLootBoxCollectionFromVortexDir(dirPath string) (map[string]*LootBox, error) {
	ret := make(map[string]*LootBox)
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dirPath, entry.Name()))
		if err != nil {
			return nil, err
		}
		lb, err := NewLootBoxFromVortexBytes(data)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", entry.Name(), err)
		}
		ret[lb.ID] = lb
	}
	return ret, nil
}

// NewLootBoxCollectionFromVortexFS reads raw vortex JSON files from an fs.FS.
func NewLootBoxCollectionFromVortexFS(fsys fs.FS) (map[string]*LootBox, error) {
	ret := make(map[string]*LootBox)
	entries, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		data, err := fs.ReadFile(fsys, entry.Name())
		if err != nil {
			return nil, err
		}
		lb, err := NewLootBoxFromVortexBytes(data)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", entry.Name(), err)
		}
		ret[lb.ID] = lb
	}
	return ret, nil
}
