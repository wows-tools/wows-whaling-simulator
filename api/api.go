package api

import (
	"errors"
	"github.com/kakwa/wows-whaling-simulator/config"
	"github.com/kakwa/wows-whaling-simulator/lootbox"
	"github.com/labstack/echo/v4"
)

var (
	ErrNoLootbox = errors.New("No Lootbox description found")
)

type API struct {
	echo              *echo.Echo
	cfg               *config.AppConfig
	lootboxCollection map[string]*lootbox.LootBox
}

func NewAPI(echo *echo.Echo, cfg *config.AppConfig) (*API, error) {
	var a API
	var err error
	a.echo = echo
	a.cfg = cfg
	a.lootboxCollection, err = lootbox.NewLootBoxCollection(a.cfg.Rates)
	if err != nil {
		return nil, err
	}
	if len(a.lootboxCollection) == 0 {
		return nil, ErrNoLootbox
	}

	return &a, nil
}
