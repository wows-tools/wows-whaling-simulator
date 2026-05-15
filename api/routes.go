package api

import (
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

func (a *API) RegisterRoutes() {
	a.echo.GET("/ping", pingHandler)
	a.echo.GET("/api/v1/stats", a.getStats)
	a.echo.GET("/api/v1/realms/:realm_id/players", a.searchPlayer)
	a.echo.GET("/api/v1/lootboxes", a.listLootboxes)
	a.echo.GET("/api/v1/lootboxes/:lootbox_id", a.getLootbox)
	a.echo.GET("/api/v1/lootboxes/:lootbox_id/collectables", a.listLootboxCollectables)
	a.echo.GET("/api/v1/lootboxes/:lootbox_id/realms/:realm_id/players/:player_id/simple_whaling_quantity", a.simpleWhalingQuantity)
	a.echo.GET("/api/v1/lootboxes/:lootbox_id/realms/:realm_id/players/:player_id/stats_whaling_quantity", a.statsWhalingQuantity)
	a.echo.GET("/api/v1/lootboxes/:lootbox_id/realms/:realm_id/players/:player_id/simple_whaling_target", a.simpleWhalingTarget)
	a.echo.GET("/api/v1/lootboxes/:lootbox_id/realms/:realm_id/players/:player_id/stats_whaling_target", a.statsWhalingTarget)
	a.echo.GET("/api/v1/lootboxes/:lootbox_id/realms/:realm_id/players/:player_id/remainingships", a.remainingShips)

	if a.cfg.StaticAssets != "" {
		a.echo.File("/favicon.ico", filepath.Join(a.cfg.StaticAssets, "favicon.ico"))
		a.echo.File("/logo512.png", filepath.Join(a.cfg.StaticAssets, "logo512.png"))
		a.echo.File("/logo192.png", filepath.Join(a.cfg.StaticAssets, "logo192.png"))
		a.echo.File("/manifest.json", filepath.Join(a.cfg.StaticAssets, "manifest.json"))
		a.echo.File("/asset-manifest.json", filepath.Join(a.cfg.StaticAssets, "asset-manifest.json"))
		a.echo.Static("/static", filepath.Join(a.cfg.StaticAssets, "static"))
		a.echo.Static("/resources", filepath.Join(a.cfg.StaticAssets, "resources"))
		a.echo.File("*", filepath.Join(a.cfg.StaticAssets, "index.html"))
	} else {
		subFS, _ := fs.Sub(a.embeddedFS, "static")
		fileServer := http.FileServer(http.FS(subFS))
		a.echo.GET("/static/*", echo.WrapHandler(fileServer))
		a.echo.GET("/resources/*", echo.WrapHandler(fileServer))
		a.echo.GET("/favicon.ico", echo.WrapHandler(fileServer))
		a.echo.GET("/logo512.png", echo.WrapHandler(fileServer))
		a.echo.GET("/logo192.png", echo.WrapHandler(fileServer))
		a.echo.GET("/manifest.json", echo.WrapHandler(fileServer))
		a.echo.GET("/asset-manifest.json", echo.WrapHandler(fileServer))
		a.echo.GET("*", func(c echo.Context) error {
			data, err := fs.ReadFile(subFS, "index.html")
			if err != nil {
				return err
			}
			return c.Blob(http.StatusOK, "text/html; charset=utf-8", data)
		})
	}
}
