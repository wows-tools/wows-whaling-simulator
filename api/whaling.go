package api

import (
	"github.com/go-redis/cache/v8"
	"github.com/kakwa/wows-whaling-simulator/wows"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"time"
)

type BaseWhaling struct {
	Realm     string `param:"realm_id"`
	PlayerID  int    `param:"player_id"`
	LootboxID string `param:"lootbox_id"`
}

type QuantityWhaling struct {
	Quantity int `query:"number_lootbox"`
	BaseWhaling
}

func (a *API) simpleWhalingQuantity(c echo.Context) error {

	var whaling QuantityWhaling
	err := c.Bind(&whaling)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	realm, err := wows.WowsRealm(whaling.Realm)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if whaling.Quantity > 10000 || whaling.Quantity < 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Please pick a number of lootbox between 0 and 10000")
	}
	shipList := make([]string, 0)
	err = a.cache.Once(&cache.Item{
		TTL:   time.Hour,
		SetNX: true,
		Key:   "wows:searchShip:" + whaling.Realm + ":" + strconv.Itoa(whaling.PlayerID),
		Value: &shipList, // destination
		Do: func(*cache.Item) (interface{}, error) {
			shipList, err := a.wowsAPI.GetPlayerShips(realm, whaling.PlayerID)
			if err != nil {
				return nil, err
			}
			return shipList, nil
		},
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	lb, ok := a.lootboxCollection[whaling.LootboxID]
	if !ok {
		return echo.NewHTTPError(http.StatusNotFound, "Lootbox id unknown")
	}
	ws, err := lb.NewWhalingSession(shipList)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	for i := 0; i < whaling.Quantity; i++ {
		err = ws.Draw()
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	}
	ws.Finalize()
	a.stats.OpenedContainers += uint64(whaling.Quantity)
	a.stats.SimpleWhalingQuantity++

	return c.JSON(http.StatusOK, ws)
}

func (a *API) statsWhalingQuantity(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "Not implemented yet")
}

func (a *API) simpleWhalingTarget(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "Not implemented yet")
}

func (a *API) statsWhalingTarget(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "Not implemented yet")
}
