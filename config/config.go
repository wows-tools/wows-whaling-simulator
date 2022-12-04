package config

import (
	"errors"
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

type AppConfig struct {
	ListenPort   string
	RedisURI     string
	WowsApiKey   string
	LogLevel     string
	StaticAssets string
}

func (a AppConfig) ConvertLogLevel() (log.Lvl, error) {
	switch a.LogLevel {
	case "DEBUG":
		return log.DEBUG, nil
	case "INFO":
		return log.INFO, nil
	case "WARN":
		return log.WARN, nil
	case "ERROR":
		return log.ERROR, nil
	case "OFF":
		return log.OFF, nil
	default:
		return log.DEBUG, errors.New("Invalid LogLevel")
	}
}

func ParseConfig(cfgfile string) (*AppConfig, error) {

	var cfg AppConfig

	viper.SetConfigName("config")                // name of config file (without extension)
	viper.SetConfigType("yaml")                  // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/wows-whaling/")         // path to look for the config file in
	viper.AddConfigPath("$HOME/.wows-whaling")        // call multiple times to add many search paths
	viper.AddConfigPath(".")                     // optionally look for config in the working directory
	if err := viper.ReadInConfig(); err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
	viper.SetEnvPrefix("WOWSWHALING")
	viper.BindEnv("ListenPort")
	viper.BindEnv("RedisURI")
	viper.BindEnv("DBURI")
	viper.BindEnv("LogLevel")
	viper.BindEnv("StaticAssets")
	viper.AutomaticEnv()
	if err := viper.Unmarshal(&cfg); err != nil {
		fmt.Printf("unable to decode into struct, %v", err)
		return nil, err
	}
	return &cfg, nil
}
