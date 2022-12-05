package api

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (a *API) listLootboxes(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "Not implemented yet")
}

func (a *API) getLootbox(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "Not implemented yet")
}

func (a *API) listLootboxCollectables(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "Not implemented yet")
}
