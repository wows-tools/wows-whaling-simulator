package lootbox

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
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

// vortexCategoryOverrides maps exact container IDs to categories for containers
// that cannot be reliably classified by prefix alone.
var vortexCategoryOverrides = map[string]string{
	// Resources & Supplies
	"epic_bonuses_container":   "Resources & Supplies",
	"flags_premium_container":  "Resources & Supplies",
	"more_signals_container":   "Resources & Supplies",
	"more_coal_container":      "Resources & Supplies",
	"supercontainer":           "Resources & Supplies",
	"steam_supercontainer":     "Resources & Supplies",
	"steam_container":          "Resources & Supplies",

	// Collection events (ship-collection themed containers)
	"always_courageous_container":         "Collection",
	"always_courageous_premium_container": "Collection",
	"dunkirk_container":                   "Collection",
	"resolute_and_rapid_container":        "Collection",

	// Seasonal
	"journey_to_the_west_container":         "Seasonal",
	"journey_to_the_west_premium_container": "Seasonal",
	"distant_voyages_premium_container":     "Seasonal",
	"black_daruma_premium_container":        "Seasonal",
	"nians_treasure_premium_container":      "Seasonal",
	"picnic_by_the_sea_container":           "Seasonal",

	// Historical Events
	"dunkirk_premium_container":            "Historical Events",
	"operation_dynamo_container":           "Historical Events",
	"operation_dynamo_premium_container":   "Historical Events",
	"finest_hour_premium_container":        "Historical Events",

	// National Navies
	"in_the_service_of_the_motherland_container":         "National Navies",
	"in_the_service_of_the_motherland_premium_container": "National Navies",
	"french_squadron_container":                          "National Navies",
	"french_squadron_premium_container":                  "National Navies",
	"resolute_and_rapid_premium_container":               "National Navies",
	"europe_container":                                   "National Navies",

	// Collaborations
	"warhammer_40000__imperium_container":                     "Collaborations",
	"warhammer_40000__chaos_container":                        "Collaborations",
	"warhammer_40000_premium_container":                       "Collaborations",
	"two_titansone_king_premiumcontainer":                     "Collaborations",
	"star_trek_premium_container":                             "Collaborations",
	"legend_of_sabaton_container":                             "Collaborations",
	"legend_of_sabaton_premium_container":                     "Collaborations",
	"the_transformers_the_movie_35th_anniversary_container":         "Collaborations",
	"the_transformers_the_movie_35th_anniversary_premium_container": "Collaborations",

	// Premium & Special Ships (IDs that start with premium_ but are not premium_ship_*)
	"premium_shipv_elite_container": "Premium & Special Ships",
}

func vortexCategory(id string) string {
	// Check specific overrides first.
	if cat, ok := vortexCategoryOverrides[id]; ok {
		return cat
	}

	switch {
	case strings.HasPrefix(id, "santas_"):
		return "Santa's Gifts"
	case strings.HasPrefix(id, "black_friday_"):
		return "Black Friday"
	case strings.HasPrefix(id, "premium_ship_"), strings.HasPrefix(id, "special_ship_"),
		strings.HasPrefix(id, "premium_shipv_"):
		return "Premium & Special Ships"
	case strings.HasPrefix(id, "halloween_"), strings.HasPrefix(id, "jack_o_lantern"):
		return "Seasonal"
	case strings.HasPrefix(id, "lunar_new_year"), strings.HasPrefix(id, "golden_week"),
		strings.HasPrefix(id, "red_daruma"), strings.HasPrefix(id, "golden_daruma"):
		return "Seasonal"
	case strings.HasPrefix(id, "christmas_"), strings.HasPrefix(id, "festive_"):
		return "Seasonal"
	case strings.Contains(id, "_anniversary"):
		return "Anniversaries"
	case strings.HasPrefix(id, "battle_of_"), strings.HasPrefix(id, "d_day"),
		strings.HasPrefix(id, "war_is_over"), strings.HasPrefix(id, "five_epochs"),
		strings.HasPrefix(id, "the_hunt_for_bismarck"), strings.HasPrefix(id, "unsinkable_sam"),
		strings.HasPrefix(id, "allied_heroes"),
		strings.HasPrefix(id, "heart_of_oak"), strings.HasPrefix(id, "eagles"),
		strings.HasPrefix(id, "legion_of_honor"), strings.HasPrefix(id, "belle_poque"),
		strings.HasPrefix(id, "three_kingdoms"), strings.HasPrefix(id, "the_age_of_sa_zhenbing"),
		strings.HasPrefix(id, "go_navy"), strings.HasPrefix(id, "navy_250"):
		return "Historical Events"
	case strings.HasSuffix(id, "_league") || strings.Contains(id, "_league_"):
		return "Leagues"
	case strings.HasPrefix(id, "azur_lane"), strings.HasPrefix(id, "arpeggio_"),
		strings.HasPrefix(id, "high_school_fleet"), strings.HasPrefix(id, "plus_ultra"),
		strings.HasPrefix(id, "cybertronian"), strings.HasPrefix(id, "logh_"),
		strings.HasPrefix(id, "space_the_final_frontier"), strings.HasPrefix(id, "ars_nova"),
		strings.HasPrefix(id, "manjuu"), strings.HasPrefix(id, "on_a_retro_wave"),
		strings.HasPrefix(id, "psplive"), strings.HasPrefix(id, "wings_of_freedom"):
		return "Collaborations"
	case strings.HasPrefix(id, "american_"), strings.HasPrefix(id, "german_"),
		strings.HasPrefix(id, "japanese_"), strings.HasPrefix(id, "italian_"),
		strings.HasPrefix(id, "royal_navy"), strings.HasPrefix(id, "soviet"),
		strings.HasPrefix(id, "commonwealth_"), strings.HasPrefix(id, "kriegsmarine"),
		strings.HasPrefix(id, "vive_la_france"), strings.HasPrefix(id, "yamamoto_"),
		strings.HasPrefix(id, "oktyabrskaya_"):
		return "National Navies"
	case strings.Contains(id, "_recruiting_station"):
		return "Recruitment"
	case strings.HasPrefix(id, "regular_supplies"), strings.HasPrefix(id, "more_economic"),
		strings.HasPrefix(id, "more_camouflages"), strings.HasPrefix(id, "small_"),
		strings.HasPrefix(id, "common_bonuses"), strings.HasPrefix(id, "uncommon_bonuses"),
		strings.HasPrefix(id, "rare_bonuses"), strings.HasPrefix(id, "tactical"),
		strings.HasPrefix(id, "volunteer_skins"), strings.HasPrefix(id, "air_supply"),
		strings.HasPrefix(id, "epic_bonuses"), strings.HasPrefix(id, "flags_"):
		return "Resources & Supplies"
	default:
		return "Special Events"
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

	// Explicit per-container prices; fall back to 1250 for any other premium
	// variant, and 0 (unknown) for everything else.
	var price float64
	switch id {
	case "santas_gift_container":
		price = 250.0
	case "santas_big_gift_container":
		price = 750.0
	case "santas_mega_gift_container":
		price = 1250.0
	case "santas_ultra_gift_container":
		price = 1850.0
	default:
		if strings.Contains(id, "premium") && !strings.HasPrefix(id, "premium_ship") {
			price = 1250.0
		}
	}

	lb := &LootBox{
		Name:               d.Title,
		Img:                vortexIconImg(d.Icons),
		ID:                 id,
		Category:           vortexCategory(id),
		Weight:             vortexGetWeight(id),
		ExchangeRateEuro:   308.947108255,
		ExchangeRateDollar: 252.780586451,
		Pity:               pity,
		Price:              price,
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
		// Sort both the category list and the transfers slice for stable JSON output.
		sort.Strings(shipCatNames)
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
