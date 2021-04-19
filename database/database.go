package database

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/jon4hz/go-binance-local-orderbook/config"
)

var (
	db = &gorm.DB{}
)

type BinanceDepthResponse struct {
	Snapshot *binance.DepthResponse
	Response *binance.WsDepthEvent
}

type BinanceFuturesDepthResponse struct {
	Snapshot *futures.DepthResponse
	Response *futures.WsDepthEvent
}

func Connect(config *config.Config) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable ", config.Database.DBServer, config.Database.DBUser, config.Database.DBPassword, config.Database.DBName, config.Database.DBPort)
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("[database][parseConfig] Error configuring the database: ", err)
	}
}

/* func InitDatabase(config *config.Config) error {
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
		log.Println("[database][init] Successfully deleted old depth cache")
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
	log.Println("[database][init] Successfully created new tables")
	return nil
} */

type bid struct {
	Price    string `gorm:"primaryKey"`
	Quantity string
}

type ask struct {
	Price    string `gorm:"primaryKey"`
	Quantity string
}

type Tabler interface {
	TableName() string
}

func (ask) TableName() string {
	return fmt.Sprintf("%s_asks", "BTCUSDT")
}

func (bid) TableName() string {
	return fmt.Sprintf("%s_bids", "BTCUSDT")
}

func createUnifiedStruct(asks interface{}, bids interface{}) ([]ask, []bid, error) {
	jAsks, err := json.Marshal(asks)
	if err != nil {
		return nil, nil, err
	}
	oAsks := []ask{}
	err = json.Unmarshal(jAsks, &oAsks)
	if err != nil {
		return nil, nil, err
	}
	jBids, err := json.Marshal(bids)
	if err != nil {
		return nil, nil, err
	}
	oBids := []bid{}
	err = json.Unmarshal(jBids, &oBids)
	if err != nil {
		return nil, nil, err
	}
	return oAsks, oBids, nil
}

func doDBInsert(sym string, asks interface{}, bids interface{}) error {
	oAsks, oBids, err := createUnifiedStruct(asks, bids)
	if err != nil {
		return err
	}

	db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"quantity"}),
	}).Create(&oAsks)

	db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"quantity"}),
	}).Create(&oBids)

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
