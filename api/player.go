package api

import (
	wgwows "github.com/IceflowRE/go-wargaming/v3/wargaming/wows"
	"github.com/go-redis/cache/v8"
	"github.com/kakwa/wows-whaling-simulator/wows"
	"github.com/labstack/echo/v4"

	"net/http"
	"strconv"
	"time"
)

type UserSearch struct {
	NickStart string `query:"nick_start"`
	Realm     string `param:"realm_id"`
}

type SearchResult struct {
	Players []*wgwows.AccountList `json:"players"`
}

type UserShipList struct {
	Realm    string `param:"realm_id"`
	PlayerID int    `param:"player_id"`
}

type ShipListResult struct {
	Ships []string `json:"ships"`
}

func (a *API) searchPlayer(c echo.Context) error {
	var user UserSearch
	err := c.Bind(&user)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	realm, err := wows.WowsRealm(user.Realm)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	accountList := make([]*wgwows.AccountList, 0)
	err = a.cache.Once(&cache.Item{
		TTL:   time.Hour,
		SetNX: true,
		Key:   "wows:searchPlayer:" + user.Realm + ":" + user.NickStart,
		Value: &accountList, // destination
		Do: func(*cache.Item) (interface{}, error) {
			accountList, err := a.wowsAPI.SearchPlayer(realm, user.NickStart, "startswith")
			if err != nil {
				return nil, err
			}
			return accountList, nil
		},
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ret := SearchResult{
		Players: accountList,
	}
	return c.JSON(http.StatusOK, ret)
}

func (a *API) listPlayerShips(c echo.Context) error {
	var user UserShipList
	err := c.Bind(&user)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	realm, err := wows.WowsRealm(user.Realm)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	shipList := make([]string, 0)
	err = a.cache.Once(&cache.Item{
		TTL:   time.Hour,
		SetNX: true,
		Key:   "wows:searchShip:" + user.Realm + ":" + strconv.Itoa(user.PlayerID),
		Value: &shipList, // destination
		Do: func(*cache.Item) (interface{}, error) {
			shipList, err := a.wowsAPI.GetPlayerShips(realm, user.PlayerID)
			if err != nil {
				return nil, err
			}
			return shipList, nil
		},
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ret := ShipListResult{
		Ships: shipList,
	}
	return c.JSON(http.StatusOK, ret)
}
