package main

import (
	"github.com/kakwa/wows-whaling-simulator/api"
	"github.com/kakwa/wows-whaling-simulator/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	// TODO properly set log level
	e.HideBanner = true

	cfg, err := config.ParseConfig("./config/example.yml")

	if err != nil {
		e.Logger.Fatalf("failed to parse configuration: %w", err)
	}
	logLevel, err := cfg.ConvertLogLevel()
	if err != nil {
		e.Logger.Fatalf("failed to parse LogLevel: %s", err)
	}

	e.Logger.SetLevel(logLevel)

	_, err = api.NewAPI(e, cfg)
	if err != nil {
		e.Logger.Fatalf("failed to init API: %s", err)
	}

	e.Logger.Fatal(e.Start(":" + cfg.ListenPort))
}
