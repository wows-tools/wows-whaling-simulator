package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/IceflowRE/go-wargaming/v3/wargaming"
	wgwows "github.com/IceflowRE/go-wargaming/v3/wargaming/wows"
)

func main() {
	collect := flag.Bool("collect", false, "scrape WG page and fetch lootbox JSON into input dir")
	inputDir := flag.String("input", "raw", "directory with raw WG JSON files")
	outputDir := flag.String("output", "output", "directory for converted JSON and images")
	ratesDir := flag.String("rates", "../../rates", "rates directory to copy Santa containers into")
	wgKey := flag.String("wg-key", "", "Wargaming API application_id (for ship name resolution)")
	flag.Parse()

	imgDir := filepath.Join(*outputDir, "resources")
	for _, dir := range []string{*inputDir, *outputDir, imgDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Println("Error creating directory:", err)
			return
		}
	}

	if *collect {
		if err := collectRaw(*inputDir); err != nil {
			fmt.Println("Collection failed:", err)
			return
		}
	}

	// Build ship name+tier map from WG API if a key was provided.
	ships := make(map[int]ShipInfo)
	if *wgKey != "" {
		fmt.Println("Fetching ship data from WG API…")
		var err error
		ships, err = fetchShips(*wgKey)
		if err != nil {
			fmt.Println("Warning: could not fetch ship data:", err)
		} else {
			fmt.Printf("Loaded %d ships\n", len(ships))
		}
	}

	entries, err := os.ReadDir(*inputDir)
	if err != nil {
		fmt.Println("Error reading input dir:", err)
		return
	}

	santaYear := santaContainerYear()

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		inPath := filepath.Join(*inputDir, entry.Name())
		raw, err := os.ReadFile(inPath)
		if err != nil {
			fmt.Printf("Error reading %s: %v\n", inPath, err)
			continue
		}

		var lb WgLootbox
		if err := json.Unmarshal(raw, &lb); err != nil {
			fmt.Printf("Error parsing %s: %v\n", inPath, err)
			continue
		}

		lbws, err := convert(&lb, ships)
		if err != nil {
			fmt.Printf("Error converting %s: %v\n", inPath, err)
			continue
		}

		// Download the lootbox icon.
		if err := downloadIcon(lb.Data.Icons, imgDir); err != nil {
			fmt.Printf("Warning: could not download icon for %s: %v\n", lb.Data.Title, err)
		}

		outPath := filepath.Join(*outputDir, lbws.ID+".json")
		data, err := json.MarshalIndent(lbws, "", "    ")
		if err != nil {
			fmt.Printf("Error marshaling %s: %v\n", inPath, err)
			continue
		}
		if err := os.WriteFile(outPath, data, 0644); err != nil {
			fmt.Printf("Error writing %s: %v\n", outPath, err)
			continue
		}
		fmt.Printf("Converted %s → %s\n", inPath, outPath)

		// Copy Santa containers to the rates directory with the year suffix.
		if ratesName := santaRatesName(lb.Data.Title, santaYear); ratesName != "" {
			ratesPath := filepath.Join(*ratesDir, ratesName)
			if err := os.WriteFile(ratesPath, data, 0644); err != nil {
				fmt.Printf("Warning: could not write rates file %s: %v\n", ratesPath, err)
			} else {
				fmt.Printf("  → rates: %s\n", ratesPath)
			}
		}
	}
}

// santaContainerYear returns the year to use for Santa container filenames.
// Christmas containers are released in December; if we're in the first half of
// the year the current containers are from the previous December.
func santaContainerYear() int {
	now := time.Now()
	if now.Month() < time.July {
		return now.Year() - 1
	}
	return now.Year()
}

// santaRatesName maps a container title to a rates/ filename (empty = not a Santa container).
func santaRatesName(title string, year int) string {
	t := strings.ToLower(title)
	switch {
	case strings.Contains(t, "santa") && strings.Contains(t, "mega"):
		return fmt.Sprintf("santa_mega_%d.json", year)
	case strings.Contains(t, "santa") && strings.Contains(t, "big"):
		return fmt.Sprintf("santa_big_%d.json", year)
	case strings.Contains(t, "santa"):
		return fmt.Sprintf("santa_%d.json", year)
	}
	return ""
}

// fetchShips queries the WG encyclopedia for all ship names, tiers, and rarity flags.
func fetchShips(key string) (map[int]ShipInfo, error) {
	client := wargaming.NewClient(key, &wargaming.ClientOptions{
		HTTPClient: &http.Client{Timeout: 15 * time.Second},
	})
	result := make(map[int]ShipInfo)
	pageNo := 1
	limit := 100
	lang := "en"
	for {
		res, err := client.Wows.EncyclopediaShips(
			context.Background(),
			wargaming.RealmEu,
			&wgwows.EncyclopediaShipsOptions{
				Fields:   []string{"ship_id", "name", "tier", "is_premium", "is_special"},
				PageNo:   &pageNo,
				Limit:    &limit,
				Language: &lang,
			},
		)
		if err != nil || len(res) == 0 {
			break
		}
		for _, s := range res {
			if s.ShipId == nil || s.Name == nil {
				continue
			}
			info := ShipInfo{Name: *s.Name}
			if s.Tier != nil {
				info.Tier = *s.Tier
			}
			if s.IsPremium != nil {
				info.IsPremium = *s.IsPremium
			}
			if s.IsSpecial != nil {
				info.IsSpecial = *s.IsSpecial
			}
			result[*s.ShipId] = info
		}
		pageNo++
	}
	return result, nil
}

// collectRaw uses Selenium to scrape the WG lootbox page, then fetches each
// discovered vortex API URL and saves the raw JSON into destDir.
func collectRaw(destDir string) error {
	fmt.Println("Launching browser to collect lootbox URLs…")
	urls := CollectLootboxURLs()
	fmt.Printf("Found %d lootbox URL(s)\n", len(urls))
	for i, u := range urls {
		fmt.Printf("Fetching [%d/%d] %s\n", i+1, len(urls), u)
		resp, err := http.Get(u) //nolint:gosec
		if err != nil {
			return fmt.Errorf("fetch %s: %w", u, err)
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return fmt.Errorf("read %s: %w", u, err)
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("fetch %s: status %d", u, resp.StatusCode)
		}
		outPath := filepath.Join(destDir, fmt.Sprintf("%d.json", i))
		if err := os.WriteFile(outPath, body, 0644); err != nil {
			return fmt.Errorf("write %s: %w", outPath, err)
		}
		fmt.Printf("  → saved %s\n", outPath)
	}
	return nil
}

const wgIconBase = "https://wows-gloss-icons.wgcdn.co/"

// downloadIcon fetches all available icon variants and saves them under imgDir.
func downloadIcon(icons Icon, imgDir string) error {
	urls := []string{icons.Default, icons.Large, icons.Small}
	seen := map[string]bool{}
	var errs []string
	for _, u := range urls {
		if u == "" || seen[u] {
			continue
		}
		seen[u] = true
		if err := downloadFile(u, imgDir); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("%s", strings.Join(errs, "; "))
	}
	return nil
}

// downloadFile fetches url and writes the response body to imgDir/<basename>.
// Relative paths are resolved against wgIconBase.
func downloadFile(url, destDir string) error {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = wgIconBase + url
	}
	base := filepath.Base(url)
	if idx := strings.Index(base, "?"); idx != -1 {
		base = base[:idx]
	}
	dest := filepath.Join(destDir, base)

	if _, err := os.Stat(dest); err == nil {
		return nil
	}

	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return fmt.Errorf("GET %s: %w", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GET %s: status %d", url, resp.StatusCode)
	}

	f, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("create %s: %w", dest, err)
	}
	defer f.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		return fmt.Errorf("write %s: %w", dest, err)
	}
	return nil
}

func name2id(input string) string {
	result := strings.ToLower(input)
	// Strip ASCII and Unicode typographic quotes.
	for _, q := range []string{"'", "’", "\"", "“", "”", "«", "»"} {
		result = strings.ReplaceAll(result, q, "")
	}
	// Collapse spaces to underscores, drop anything else non-word.
	var b strings.Builder
	for _, r := range result {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == ' ' || r == '_' || r == '-':
			b.WriteByte('_')
		}
	}
	return b.String()
}

func getWeight(id string) int {
	if strings.Contains(id, "mega") {
		return 40
	}
	if strings.Contains(id, "big") {
		return 110
	}
	return 1100
}
