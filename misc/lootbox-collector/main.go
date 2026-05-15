package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/IceflowRE/go-wargaming/v3/wargaming"
	wgwows "github.com/IceflowRE/go-wargaming/v3/wargaming/wows"
	"github.com/kakwa/wows-whaling-simulator/lootbox"
)

// ShipEntry is the per-ship record written to ships.json.
type ShipEntry struct {
	Name      string `json:"name"`
	Tier      int    `json:"tier"`
	IsPremium bool   `json:"is_premium"`
	IsSpecial bool   `json:"is_special"`
}

func main() {
	collect := flag.Bool("collect", false, "scrape WG page and fetch lootbox JSON into input dir")
	inputDir := flag.String("input", "raw", "directory with raw WG JSON files")
	outputDir := flag.String("output", "output", "directory for converted JSON")
	ratesDir := flag.String("rates", "../../rates", "rates directory to copy all containers into")
	staticDir := flag.String("static", "../../static", "static assets directory; icons are written to <static>/resources/")
	shipsFile := flag.String("ships-file", "../../ships/ships.json", "output file for ship ID→name/tier mapping")
	wgKey := flag.String("wg-key", "", "Wargaming API application_id (for ship name resolution)")
	flag.Parse()

	imgDir := filepath.Join(*staticDir, "resources")
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

	// Pre-scan: find the last raw filename for each container ID so that when
	// duplicate raw files exist (old + new version of the same container) only
	// the last one triggers change-detection in rates/.
	lastFileForID := map[string]string{}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(*inputDir, entry.Name()))
		if err != nil {
			continue
		}
		var env struct {
			Data struct{ Title string `json:"title"` } `json:"data"`
		}
		if json.Unmarshal(raw, &env) == nil && env.Data.Title != "" {
			id := vortexName2IDPublic(env.Data.Title)
			lastFileForID[id] = entry.Name()
		}
	}

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

		lbws, err := convert(raw, ships)
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

		// Copy to rates only from the last raw file for this ID so duplicate
		// raw files don't cause perpetual false change-detections.
		if *ratesDir != "" && lastFileForID[lbws.ID] == entry.Name() {
			if err := writeToRates(*ratesDir, lbws, data); err != nil {
				fmt.Printf("Warning: could not write rates file %s: %v\n", lbws.ID, err)
			}
		}
	}

	if *shipsFile != "" {
		rawShips := extractShipsFromRaw(*inputDir)
		if err := writeShipsFile(*shipsFile, rawShips, ships); err != nil {
			fmt.Printf("Warning: could not write ships file: %v\n", err)
		} else {
			total := len(rawShips)
			if len(ships) > total {
				total = len(ships)
			}
			fmt.Printf("Wrote %d ships to %s\n", total, *shipsFile)
		}
	}
}

// extractShipsFromRaw scans all raw vortex JSON files in inputDir and builds a
// map of ship ID → ShipEntry from the additionalData embedded in each file.
func extractShipsFromRaw(inputDir string) map[string]ShipEntry {
	ships := make(map[string]ShipEntry)
	entries, err := os.ReadDir(inputDir)
	if err != nil {
		return ships
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(inputDir, entry.Name()))
		if err != nil {
			continue
		}
		var lb WgLootbox
		if err := json.Unmarshal(raw, &lb); err != nil {
			continue
		}
		for _, slot := range lb.Data.Slots {
			for _, vlist := range slot.ValuableRewards {
				for _, r := range vlist.Rewards {
					if r.Type != "ship" || r.AdditionalData == nil || r.AdditionalData.Title == "" {
						continue
					}
					ships[strconv.Itoa(r.ID)] = ShipEntry{
						Name:      r.AdditionalData.Title,
						Tier:      r.AdditionalData.Level,
						IsPremium: r.AdditionalData.IsPremium,
						IsSpecial: r.AdditionalData.IsSpecial,
					}
				}
			}
		}
	}
	return ships
}

// writeShipsFile merges ship data extracted from raw files with optional WG API
// results and writes the combined map to path as JSON.
func writeShipsFile(path string, fromRaw map[string]ShipEntry, fromAPI map[int]ShipInfo) error {
	merged := make(map[string]ShipEntry, len(fromRaw)+len(fromAPI))
	for k, v := range fromRaw {
		merged[k] = v
	}
	// API data takes precedence (more authoritative names/tiers).
	for id, info := range fromAPI {
		merged[strconv.Itoa(id)] = ShipEntry{
			Name:      info.Name,
			Tier:      info.Tier,
			IsPremium: info.IsPremium,
			IsSpecial: info.IsSpecial,
		}
	}
	data, err := json.MarshalIndent(merged, "", "    ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// writeToRates writes the container JSON to <ratesDir>/<id>.json.
// If a file with the same ID already exists but has different LootBox content,
// the old file is archived as <id>_preYYYYMMDD.json and the new file includes
// "changed_at" and "archived_previous_as" metadata fields.
func writeToRates(ratesDir string, lb *lootbox.LootBox, data []byte) error {
	if err := os.MkdirAll(ratesDir, 0755); err != nil {
		return err
	}
	dest := filepath.Join(ratesDir, lb.ID+".json")

	if existing, err := os.ReadFile(dest); err == nil {
		var existingLB lootbox.LootBox
		if jerr := json.Unmarshal(existing, &existingLB); jerr == nil {
			existBytes, _ := json.Marshal(existingLB)
			newBytes, _ := json.Marshal(*lb)
			if !bytes.Equal(existBytes, newBytes) {
				today := time.Now().Format("20060102")
				archiveName := lb.ID + "_pre" + today + ".json"
				archivePath := filepath.Join(ratesDir, archiveName)
				if rerr := os.Rename(dest, archivePath); rerr != nil {
					return fmt.Errorf("archive %s: %w", dest, rerr)
				}
				fmt.Printf("  → rates: changed, archived previous → %s\n", archiveName)
				data = injectChangeMeta(data, today, archiveName)
			}
		}
	}

	if werr := os.WriteFile(dest, data, 0644); werr != nil {
		return werr
	}
	fmt.Printf("  → rates: %s\n", dest)
	return nil
}

// vortexName2IDPublic mirrors the lootbox package's internal name→ID helper so
// the pre-scan can predict container IDs without a full parse.
func vortexName2IDPublic(s string) string {
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

// injectChangeMeta adds change-tracking fields to a JSON blob without touching
// the rest of the structure.
func injectChangeMeta(data []byte, today, archiveName string) []byte {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(data, &obj); err != nil {
		return data
	}
	obj["changed_at"] = json.RawMessage(`"` + today + `"`)
	obj["archived_previous_as"] = json.RawMessage(`"` + archiveName + `"`)
	out, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		return data
	}
	return out
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

// sessionCookies holds cookies extracted from the browser session for
// authenticated icon downloads. Set during -collect, empty otherwise.
var sessionCookies []*http.Cookie

// collectRaw uses Selenium to scrape the WG lootbox page, then fetches each
// discovered vortex API URL and saves the raw JSON into destDir.
func collectRaw(destDir string) error {
	fmt.Println("Launching browser to collect lootbox URLs…")
	urls, cookies := CollectLootboxURLs()
	sessionCookies = cookies
	fmt.Printf("Found %d lootbox URL(s), %d session cookies\n", len(urls), len(cookies))
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

const wgIconBase = "https://wows-gloss-icons.wgcdn.co/icons/"

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

	req, err := http.NewRequest(http.MethodGet, url, nil) //nolint:gosec
	if err != nil {
		return fmt.Errorf("GET %s: %w", url, err)
	}
	req.Header.Set("Referer", "https://worldofwarships.com/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	for _, c := range sessionCookies {
		req.AddCookie(c)
	}
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
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
