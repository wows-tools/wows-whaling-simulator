package api

import (
	"errors"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/kakwa/wows-whaling-simulator/config"
	"github.com/kakwa/wows-whaling-simulator/lootbox"
	"github.com/kakwa/wows-whaling-simulator/wows"
	"github.com/labstack/echo/v4"
	"time"
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
	lootboxCollection map[string]*lootbox.LootBox
	redis             *redis.Client
	cache             *cache.Cache
	wowsAPI           *wows.WowsAPI
	stats             *Stats
}

func NewAPI(echo *echo.Echo, cfg *config.AppConfig) (*API, error) {
	var a API
	var err error
	a.echo = echo
	a.cfg = cfg
	a.stats = &Stats{}
	a.lootboxCollection, err = lootbox.NewLootBoxCollection(a.cfg.Rates)
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
