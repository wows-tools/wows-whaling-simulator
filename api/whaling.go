package api

import (
	mapset "github.com/deckarep/golang-set/v2"
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

type TargetWhaling struct {
	Target        string   `query:"target"`
	ExcludedShips []string `query:"excluded_ships[]"`
	BaseWhaling
}

type RemainingShips struct {
	Ships []string `json:"ships"`
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

	ws.SimpleWhaling(whaling.Quantity)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	a.stats.OpenedContainers += uint64(whaling.Quantity)
	a.stats.WhalingQuantitySimple++

	return c.JSON(http.StatusOK, ws)
}

func (a *API) statsWhalingQuantity(c echo.Context) error {
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
	wss := lb.NewWhalingStatsSession(shipList)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	wss.StatsSimpleWhaling(whaling.Quantity)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	a.stats.OpenedContainers += wss.Opened
	a.stats.WhalingQuantityStats++

	return c.JSON(http.StatusOK, wss)
}

func (a *API) simpleWhalingTarget(c echo.Context) error {
	var whaling TargetWhaling
	err := c.Bind(&whaling)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	realm, err := wows.WowsRealm(whaling.Realm)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if len(whaling.Target) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing target for this whaling")
	}
	shipList := make([]string, 0)
	err = a.cache.Once(&cache.Item{
		TTL:   time.Minute * 10,
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
	setApiShipList := mapset.NewSet[string](shipList...)
	setExcludeShipList := mapset.NewSet[string](whaling.ExcludedShips...)
	setUnionExclusion := setApiShipList.Union(setExcludeShipList)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	lb, ok := a.lootboxCollection[whaling.LootboxID]
	if !ok {
		return echo.NewHTTPError(http.StatusNotFound, "Lootbox id unknown")
	}
	ws, err := lb.NewWhalingSession(setUnionExclusion.ToSlice())
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = ws.TargetWhaling(whaling.Target)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	a.stats.OpenedContainers += ws.ContainerOpened
	a.stats.TargetWhalingSimple++

	return c.JSON(http.StatusOK, ws)
}

func (a *API) statsWhalingTarget(c echo.Context) error {
	var whaling TargetWhaling
	err := c.Bind(&whaling)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	realm, err := wows.WowsRealm(whaling.Realm)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if len(whaling.Target) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing target for this whaling")
	}
	shipList := make([]string, 0)
	err = a.cache.Once(&cache.Item{
		TTL:   time.Minute * 10,
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
	setApiShipList := mapset.NewSet[string](shipList...)
	setExcludeShipList := mapset.NewSet[string](whaling.ExcludedShips...)
	setUnionExclusion := setApiShipList.Union(setExcludeShipList)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	lb, ok := a.lootboxCollection[whaling.LootboxID]
	if !ok {
		return echo.NewHTTPError(http.StatusNotFound, "Lootbox id unknown")
	}
	wss := lb.NewWhalingStatsSession(setUnionExclusion.ToSlice())
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = wss.StatsTargetWhaling(whaling.Target)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	a.stats.OpenedContainers += wss.Opened
	a.stats.TargetWhalingStats++

	return c.JSON(http.StatusOK, wss)

}

func (a *API) remainingShips(c echo.Context) error {

	var whaling BaseWhaling
	err := c.Bind(&whaling)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	realm, err := wows.WowsRealm(whaling.Realm)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
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

	ret := &RemainingShips{
		Ships: ws.RemainingCollectables(),
	}
	return c.JSON(http.StatusOK, ret)
}
