package main

import (
	"os"
	"strings"

	"github.com/jon4hz/go-binance-local-orderbook/config"
	"github.com/jon4hz/go-binance-local-orderbook/database"
	"github.com/jon4hz/go-binance-local-orderbook/exchange/binance"
	"github.com/jon4hz/go-binance-local-orderbook/watchdog"
)

func main() {
	// load configuration from file or env
	config := loadConfiguration()

	// get database pool
	database.Connect(config)
	/* if err != nil {
		os.Exit(1)
	} */

	/* err = database.InitDatabase(config)
	if err != nil {
		log.Fatal(err)
	} */
	// start orderbook watchdog
	go watchdog.Watcher(config)

	ch := make(chan bool)
	// change binance to exchange package
	go binance.InitWebsocket(config)
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
