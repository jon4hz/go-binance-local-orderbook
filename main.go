package main

import (
	"os"

	"github.com/jon4hz/go-binance-local-orderbook/config"
	"github.com/jon4hz/go-binance-local-orderbook/exchange/binance"
	"github.com/jon4hz/go-binance-local-orderbook/watchdog"
)

func main() {
	config := loadConfiguration()
	go watchdog.Watcher()
	binance.HandleWebsocket(config)
}

func loadConfiguration() *config.Config {
	var err error
	customConfigFile := os.Getenv("CONFIG_FILE")
	if len(customConfigFile) > 0 {
		err = config.Load(customConfigFile)
	} else {
		err = config.Load(config.DefaultConfigurationFilePath)
	}
	if err != nil {
		panic(err)
	}
	return config.Get()
}
