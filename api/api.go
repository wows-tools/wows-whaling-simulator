package api

import (
	"errors"
	"io/fs"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/kakwa/wows-whaling-simulator/config"
	"github.com/kakwa/wows-whaling-simulator/lootbox"
	"github.com/kakwa/wows-whaling-simulator/wows"
	"github.com/labstack/echo/v4"
)

var (
	ErrNoLootbox = errors.New("No Lootbox description found")
)

type Stats struct {
	WhalingQuantitySimple uint64 `json:"simple_whaling_quantity"`
	WhalingQuantityStats  uint64 `json:"stats_whaling_quantity"`
	TargetWhalingSimple   uint64 `json:"simple_whaling_target"`
	TargetWhalingStats    uint64 `json:"stats_whaling_target"`
	OpenedContainers      uint64 `json:"opened_counters"`
}

type API struct {
	echo              *echo.Echo
	cfg               *config.AppConfig
	embeddedFS        fs.FS
	lootboxCollection map[string]*lootbox.LootBox
	redis             *redis.Client
	cache             *cache.Cache
	wowsAPI           *wows.WowsAPI
	stats             *Stats
}

func NewAPI(echo *echo.Echo, cfg *config.AppConfig, embeddedFS fs.FS) (*API, error) {
	var a API
	var err error
	a.echo = echo
	a.cfg = cfg
	a.embeddedFS = embeddedFS
	a.stats = &Stats{}

	if cfg.Rates != "" {
		a.lootboxCollection, err = lootbox.NewLootBoxCollection(cfg.Rates)
	} else {
		ratesFS, fsErr := fs.Sub(embeddedFS, "rates")
		if fsErr != nil {
			return nil, fsErr
		}
		a.lootboxCollection, err = lootbox.NewLootBoxCollectionFromFS(ratesFS)
	}
	if err != nil {
		return nil, err
	}
	if len(a.lootboxCollection) == 0 {
		return nil, ErrNoLootbox
	}
	opt, err := redis.ParseURL(a.cfg.RedisURI)
	if err != nil {
		return nil, err
	}

	a.redis = redis.NewClient(opt)
	a.cache = cache.New(&cache.Options{
		Redis:      a.redis,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})
	a.wowsAPI = wows.NewWowsAPI(a.cfg.WowsApiKey)
	// Pre-populate ship mapping from the embedded ships.json so that ship names
	// are available even when the WG API is unreachable or the key is missing.
	if loadErr := a.wowsAPI.LoadShipMappingFromFS(embeddedFS, "ships/ships.json"); loadErr != nil {
		// Non-fatal: the API call below will still populate what it can.
		_ = loadErr
	}
	shipMapping := make(map[int]string)
	err = a.cache.Once(&cache.Item{
		TTL:   time.Hour,
		SetNX: true,
		Key:   "wows:shipmapping",
		Value: &shipMapping, // destination
		Do: func(*cache.Item) (interface{}, error) {
			err := a.wowsAPI.FillShipMapping()
			if err != nil {
				return nil, err
			}
			return a.wowsAPI.ShipMapping, nil
		},
	})
	if err != nil {
		return nil, err
	}

	a.wowsAPI.ShipMapping = shipMapping
	return &a, nil
}
