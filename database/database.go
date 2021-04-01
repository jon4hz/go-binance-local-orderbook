package database

import (
	"context"
	"fmt"
	"log"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jon4hz/go-binance-local-orderbook/config"
)

var (
	dbpool *pgxpool.Pool
	ctx    context.Context
)

func CreateDatabasePool(config *config.Config) (err error) {
	// create database connection
	pgxConfig, err := pgxpool.ParseConfig(fmt.Sprintf("postgres://%s:%s@%s:%s/%s", config.Database.DBUser, config.Database.DBPassword, config.Database.DBServer, config.Database.DBPort, config.Database.DBName))
	if err != nil {
		log.Fatal("Error configuring the database: ", err)
	}
	// create connection pool
	ctx = context.Background()
	dbpool, err = pgxpool.ConnectConfig(ctx, pgxConfig)
	if err != nil {
		log.Fatal("Unable to connect to database: ", err)
		return
	}
	return
}

func InitDatabase(config *config.Config) error {
	// drop old tables if set to true (default)
	conn, err := dbpool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()
	if config.Database.DBDeleteOldSnap {
		tables := [3]string{"asks", "bids", "general"}
		for _, table := range tables {
			if _, err := conn.Exec(ctx,
				fmt.Sprintf(`DROP TABLE IF EXISTS %s_%s;`, config.Database.DBTableMarketName, table)); err != nil {
				return err
			}
		}
		log.Println("Successfully deleted old depth cache")
	}

	// create tables
	if _, err := conn.Exec(ctx,
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s_general(
			id integer NOT NULL DEFAULT '1',
			lastUpdateId bigint NOT NULL,
			u_small bigint DEFAULT '0',
			u_big bigint DEFAULT '0',
			PRIMARY KEY(id)
		);
		
		CREATE TABLE IF NOT EXISTS %s_asks(
			id float,
			value float,
			PRIMARY KEY(id)
		);
		
		CREATE TABLE IF NOT EXISTS %s_bids(
			id float,
			value float,
			PRIMARY KEY(id)
		);`, config.Database.DBTableMarketName, config.Database.DBTableMarketName, config.Database.DBTableMarketName)); err != nil {
		return err
	}
	log.Println("Successfully created new tables")
	return nil
}

type DatabaseInsert interface {
	InsertIntoDatabase()
}

type BinanceDepthResponse struct {
	Snapshot *binance.DepthResponse
	Response *binance.WsDepthEvent
}

type BinanceFuturesResponse struct {
	Snapshot *futures.DepthResponse
	Response *futures.WsDepthEvent
}

func (data *BinanceDepthResponse) InsertIntoDatabase() {
	for _, i := range data.Response.Asks {
		fmt.Println(i.Price)
	}
}

func (data *BinanceFuturesResponse) InsertIntoDatabase() {

}
