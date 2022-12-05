package api

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (a *API)simpleWhalingQuantity(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "Not implemented yet")
}

func (a *API)statsWhalingQuantity(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "Not implemented yet")
}

func (a *API)simpleWhalingTarget(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "Not implemented yet")
}

func (a *API)statsWhalingTarget(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "Not implemented yet")
}
