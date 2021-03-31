package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/jon4hz/go-binance-local-orderbook/config"
	"github.com/jon4hz/go-binance-local-orderbook/database"
	"github.com/jon4hz/go-binance-local-orderbook/exchange/binance"
	"github.com/jon4hz/go-binance-local-orderbook/watchdog"
)

func main() {
	// load configuration from file or env
	config := loadConfiguration()

	// get database pool
	dbpool, err := createDatabasePool(config)
	if err != nil {
		os.Exit(1)
	}

	err = database.InitDatabase(dbpool)

	// start orderbook watchdog
	go watchdog.Watcher()

	// start the websocket to binance (blocking with channel)
	binance.HandleWebsocket(config)
}

func createDatabasePool(config config.Config) (dbpool *pgxpool.Pool, err error) {
	// create database connection
	pgxConfig, err := pgxpool.ParseConfig(fmt.Sprintf("postgres://%s:%s@%s:%s/%s", config.Database.DBUser, config.Database.DBPassword, config.Database.DBServer, config.Database.DBPort, config.Database.DBName))
	if err != nil {
		log.Fatal("Error configuring the database: ", err)
	}
	// create connection pool
	dbpool, err = pgxpool.ConnectConfig(context.Background(), pgxConfig)
	if err != nil {
		log.Fatal("Unable to connect to database: ", err)
		return
	}
	return
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
