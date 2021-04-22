package main

import (
	"log"
	"os"
	"strings"

	"github.com/jon4hz/go-binance-local-orderbook/config"
	"github.com/jon4hz/go-binance-local-orderbook/database"
	"github.com/jon4hz/go-binance-local-orderbook/exchange/binance"
	"github.com/jon4hz/go-binance-local-orderbook/exchange/binance_futures"
	"github.com/jon4hz/go-binance-local-orderbook/watchdog"
)

func main() {
	// load configuration from file or env
	config := loadConfiguration()

	// get database pool
	database.Connect(config.Database)

	err := database.Init(config.Database)
	if err != nil {
		log.Fatal(err)
	}
	// start orderbook watchdog
	go watchdog.Watcher(config)

	ch := make(chan bool)
	startSocketInBackground(config)
	<-ch
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
	cfg := config.Get()
	cfg.Database.DBTableMarketName = strings.ToLower(cfg.Exchange.Market)
	cfg.Database.DBDeleteOldSnap = cfg.DeleteOldSnap
	return cfg
}

func startSocketInBackground(cfg *config.Config) {
	switch cfg.Exchange.Name {
	case "binance-futures":
		go binance_futures.InitWebsocket(cfg)
	default:
		go binance.InitWebsocket(cfg)
	}

}
