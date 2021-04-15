package database

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jon4hz/go-binance-local-orderbook/config"
)

var (
	dbpool *pgxpool.Pool
)

type BinanceDepthResponse struct {
	Snapshot *binance.DepthResponse
	Response *binance.WsDepthEvent
}

type BinanceFuturesDepthResponse struct {
	Snapshot *futures.DepthResponse
	Response *futures.WsDepthEvent
}

func CreateDatabasePool(config *config.Config) (err error) {
	// create database connection
	pgxConfig, err := pgxpool.ParseConfig(fmt.Sprintf("postgres://%s:%s@%s:%s/%s", config.Database.DBUser, config.Database.DBPassword, config.Database.DBServer, config.Database.DBPort, config.Database.DBName))
	if err != nil {
		log.Fatal("Error configuring the database: ", err)
	}
	// create connection pool
	ctx := context.Background()
	dbpool, err = pgxpool.ConnectConfig(ctx, pgxConfig)
	if err != nil {
		log.Fatal("Unable to connect to database: ", err)
		return
	}
	return
}

func InitDatabase(config *config.Config) error {
	// drop old tables if set to true (default)
	ctx := context.TODO()
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
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s_asks(
			id varchar(50),
			quantity varchar(50),
			PRIMARY KEY(id)
		);
		
		CREATE TABLE IF NOT EXISTS %s_bids(
			id varchar(50),
			quantity varchar(50),
			PRIMARY KEY(id)
		);`, config.Database.DBTableMarketName, config.Database.DBTableMarketName)); err != nil {
		return err
	}
	log.Println("Successfully created new tables")
	return nil
}

type bid struct {
	Price    string
	Quantity string
}

type ask struct {
	Price    string
	Quantity string
}

func createUnifiedStruct(asks interface{}, bids interface{}) ([]*ask, []*bid, error) {
	jAsks, err := json.Marshal(asks)
	if err != nil {
		return nil, nil, err
	}
	oAsks := []*ask{}
	err = json.Unmarshal(jAsks, &oAsks)
	if err != nil {
		return nil, nil, err
	}
	jBids, err := json.Marshal(bids)
	if err != nil {
		return nil, nil, err
	}
	oBids := []*bid{}
	err = json.Unmarshal(jBids, &oBids)
	if err != nil {
		return nil, nil, err
	}
	return oAsks, oBids, nil
}

func doDBInsert(sym string, asks interface{}, bids interface{}) error {
	conn, err := dbpool.Acquire(context.TODO())
	if err != nil {
		return err
	}
	defer conn.Release()
	// create unified structs
	oAsks, oBids, err := createUnifiedStruct(asks, bids)
	if err != nil {
		return err
	}
	for _, v := range oAsks {
		var quant float64
		if quant, err = strconv.ParseFloat(v.Quantity, 64); err != nil {
			fmt.Printf("[database][dbinsert] couldn't convert \"quantity\" to float: %s\n", err)
			return err
		}
		if quant == 0 {
			if _, err := conn.Exec(context.TODO(),
				fmt.Sprintf("DELETE FROM %s_asks WHERE id = $1", sym), v.Price); err != nil {
				log.Printf("Error: %s", err)
				return err
			}
		} else {
			if _, err := conn.Exec(context.TODO(),
				fmt.Sprintf("INSERT INTO %s_asks(id, quantity) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET quantity = $3", sym), v.Price, v.Quantity, v.Quantity); err != nil {
				log.Printf("Error: %s", err)
				return err
			}
		}
	}
	for _, v := range oBids {
		var quant float64
		if quant, err = strconv.ParseFloat(v.Quantity, 64); err != nil {
			fmt.Printf("[database][dbinsert] couldn't convert \"quantity\" to float: %s\n", err)
			return err
		}
		if quant == 0 {
			if _, err := conn.Exec(context.TODO(),
				fmt.Sprintf("DELETE FROM %s_bids WHERE id = $1", sym), v.Price); err != nil {
				log.Printf("Error: %s", err)
				return err
			}
		} else {
			if _, err := conn.Exec(context.TODO(),
				fmt.Sprintf("INSERT INTO %s_bids(id, quantity) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET quantity = $3", sym), v.Price, v.Quantity, v.Quantity); err != nil {
				log.Printf("Error: %s", err)
				return err
			}
		}
	}
	return nil
}

func (resp *BinanceDepthResponse) InsertIntoDatabase(sym string) error {
	if resp.Snapshot != nil {
		err := doDBInsert(sym, resp.Snapshot.Asks, resp.Snapshot.Asks)
		if err != nil {
			return err
		}
	}
	if resp.Response != nil {
		err := doDBInsert(sym, resp.Response.Asks, resp.Response.Asks)
		if err != nil {
			return err
		}
	}
	return nil
}

func (resp *BinanceFuturesDepthResponse) InsertIntoDatabase(sym string) error {
	if resp.Snapshot != nil {
		err := doDBInsert(sym, resp.Snapshot.Asks, resp.Snapshot.Asks)
		if err != nil {
			return err
		}
	}
	if resp.Response != nil {
		err := doDBInsert(sym, resp.Response.Asks, resp.Response.Asks)
		if err != nil {
			return err
		}
	}
	return nil
}

func (resp *BinanceDepthResponse) DeleteFromDatabase(sym string) error {
	if resp.Snapshot != nil {
		err := doDBInsert(sym, resp.Snapshot.Asks, resp.Snapshot.Asks)
		if err != nil {
			return err
		}
	}
	if resp.Response != nil {
		err := doDBInsert(sym, resp.Response.Asks, resp.Response.Asks)
		if err != nil {
			return err
		}
	}
	return nil
}

func (resp *BinanceFuturesDepthResponse) DeleteFromDatabase(sym string) error {
	if resp.Snapshot != nil {
		err := doDBInsert(sym, resp.Snapshot.Asks, resp.Snapshot.Asks)
		if err != nil {
			return err
		}
	}
	if resp.Response != nil {
		err := doDBInsert(sym, resp.Response.Asks, resp.Response.Asks)
		if err != nil {
			return err
		}
	}
	return nil
}
