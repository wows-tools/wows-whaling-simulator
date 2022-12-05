package api

import (
	"path/filepath"
)

func (a *API) RegisterRoutes() {
	a.echo.GET("/ping", pingHandler)
	a.echo.GET("/api/v1/players", a.searchPlayer)
	a.echo.GET("/api/v1/players/:player_id/ships", a.listPlayerShips)
	a.echo.GET("/api/v1/lootboxes", a.listLootboxes)
	a.echo.GET("/api/v1/lootboxes/:lootbox_id", a.getLootbox)
	a.echo.GET("/api/v1/lootboxes/:lootbox_id/collectables", a.listLootboxCollectables)
	a.echo.POST("/api/v1/lootboxes/:lootbox_id/simple_whaling_quantity", a.simpleWhalingQuantity)
	a.echo.POST("/api/v1/lootboxes/:lootbox_id/stats_whaling_quantify", a.statsWhalingQuantity)
	a.echo.POST("/api/v1/lootboxes/:lootbox_id/simple_whaling_target", a.simpleWhalingTarget)
	a.echo.POST("/api/v1/lootboxes/:lootbox_id/stats_whaling_target", a.statsWhalingTarget)
	a.echo.File("/favicon.ico", filepath.Join(a.cfg.StaticAssets, "favicon.ico"))
	a.echo.Static("/static", filepath.Join(a.cfg.StaticAssets, "static"))
	a.echo.File("/", filepath.Join(a.cfg.StaticAssets, "index.html"))
}
