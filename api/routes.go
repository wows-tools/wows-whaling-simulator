package api

import (
	"path/filepath"
)

func (a *API) registerRoutes() {
	a.echo.GET("/ping", pingHandler)
	a.echo.File("/favicon.ico", filepath.Join(a.cfg.StaticAssets, "favicon.ico"))
	a.echo.Static("/static", filepath.Join(a.cfg.StaticAssets, "static"))
	a.echo.File("/", filepath.Join(a.cfg.StaticAssets, "index.html"))
}
