package api

import (
	"github.com/kakwa/wows-whaling-simulator/config"
	"github.com/labstack/echo/v4"
)

type API struct {
	echo *echo.Echo
	cfg  *config.AppConfig
}

func NewAPI(echo *echo.Echo, cfg *config.AppConfig) (*API, error) {
	var a API
	a.echo = echo
	a.cfg = cfg
	a.registerRoutes()
	return &a, nil

}
