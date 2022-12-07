package api

import (
	"github.com/kakwa/wows-whaling-simulator/lootbox"
	"github.com/labstack/echo/v4"

	"net/http"
)

type lootboxShort struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Img  string `json:"img"`
}

type lootboxCollection struct {
	Lootboxes []*lootboxShort `json:"lootboxes"`
}

type itemCollection struct {
	Items []*lootbox.ItemShort `json:"items"`
}

func (a *API) listLootboxes(c echo.Context) error {
	var ret lootboxCollection
	ret.Lootboxes = make([]*lootboxShort, 0)
	for _, lbfull := range a.lootboxCollection {
		lb := lootboxShort{
			Name: lbfull.Name,
			ID:   lbfull.ID,
			Img:   lbfull.Img,
		}
		ret.Lootboxes = append(ret.Lootboxes, &lb)
	}
	return c.JSON(http.StatusOK, ret)
}

func (a *API) getLootbox(c echo.Context) error {
	lbID := c.Param("lootbox_id")
	lb, ok := a.lootboxCollection[lbID]
	if !ok {
		return echo.NewHTTPError(http.StatusNotFound, "Lootbox id unknown")
	}
	return c.JSON(http.StatusOK, lb)
}

func (a *API) listLootboxCollectables(c echo.Context) error {
	lbID := c.Param("lootbox_id")
	lb, ok := a.lootboxCollection[lbID]
	if !ok {
		return echo.NewHTTPError(http.StatusNotFound, "Lootbox id unknown")
	}
	ret := itemCollection{
		Items: lb.ListCollectables(),
	}
	return c.JSON(http.StatusOK, ret)
}
