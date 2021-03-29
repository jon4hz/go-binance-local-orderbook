package main

import (
	"os"

	"github.com/jon4hz/go-binance-local-orderbook/config"
	"github.com/jon4hz/go-binance-local-orderbook/handler"
)

func main() {
	config := loadConfiguration()
	handler.HandleWebsocket(config)
}

func loadConfiguration() config.Config {
	var err error
	var cfg config.Config
	customConfigFile := os.Getenv("CONFIG_FILE")
	if len(customConfigFile) > 0 {
		cfg, err = config.Load(customConfigFile)
	} else {
		cfg, err = config.Load(config.DefaultConfigurationFilePath)
	}
	if err != nil {
		panic(err)
	}
	return cfg
}
