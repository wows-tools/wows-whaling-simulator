package api

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func pingHandler(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}

func (a *API) getStats(c echo.Context) error {
	return c.JSON(http.StatusOK, a.stats)
}
