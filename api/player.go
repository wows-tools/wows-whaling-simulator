package api

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (a *API)searchPlayer(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "Not implemented yet")
}

func (a *API)listPlayerShips(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "Not implemented yet")
}
