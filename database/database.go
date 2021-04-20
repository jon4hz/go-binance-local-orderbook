package database

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Config struct {
	DBName            string `mapstructure:"POSTGRES_DB"`
	DBUser            string `mapstructure:"POSTGRES_USER"`
	DBPassword        string `mapstructure:"POSTGRES_PASSWORD"`
	DBServer          string `mapstructure:"POSTGRES_SERVER"`
	DBPort            string `mapstructure:"POSTGRES_PORT"`
	DBTableMarketName string
	DBDeleteOldSnap   bool
	Debug             bool
}

var (
	db                = &gorm.DB{}
	DBTableMarketName string
)

type BinanceDepthResponse struct {
	Snapshot *binance.DepthResponse
	Response *binance.WsDepthEvent
}

type BinanceFuturesDepthResponse struct {
	Snapshot *futures.DepthResponse
	Response *futures.WsDepthEvent
}

func Connect(config *Config) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable ", config.DBServer, config.DBUser, config.DBPassword, config.DBName, config.DBPort)
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("[database][parseConfig] Error configuring the database: ", err)
	}
}

func Init(cfg *Config) (err error) {
	DBTableMarketName = cfg.DBTableMarketName

	if cfg.DBDeleteOldSnap {
		err = db.Migrator().DropTable(&ask{})
		if err != nil {
			return
		}
		err = db.Migrator().DropTable(&bid{})
		if err != nil {
			return
		}
	}
	err = db.AutoMigrate(&ask{}, &bid{})
	return
}

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
	return fmt.Sprintf("%s_asks", DBTableMarketName)
}

func (bid) TableName() string {
	return fmt.Sprintf("%s_bids", DBTableMarketName)
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

	if err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "price"}},
		DoUpdates: clause.AssignmentColumns([]string{"quantity"}),
	}).Create(&oAsks).Error; err != nil {
		return err
	}

	if err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "price"}},
		DoUpdates: clause.AssignmentColumns([]string{"quantity"}),
	}).Create(&oBids).Error; err != nil {
		return err
	}

	// loop again over bids and asks to delete 0 values
	var delAsks = []string{}
	for _, v := range oAsks {
		var quant float64
		if quant, err = strconv.ParseFloat(v.Quantity, 64); err != nil {
			log.Printf("[database][dbinsert] couldn't convert \"quantity\" to float: %s\n", err)
			return err
		}
		if quant == 0 {
			delAsks = append(delAsks, v.Price)
		}
	}
	if err := db.Delete(&oAsks, delAsks).Error; err != nil {
		return err
	}

	var delBids = []string{}
	for _, v := range oBids {
		var quant float64
		if quant, err = strconv.ParseFloat(v.Quantity, 64); err != nil {
			log.Printf("[database][dbinsert] couldn't convert \"quantity\" to float: %s\n", err)
			return err
		}
		if quant == 0 {
			delBids = append(delBids, v.Price)
		}
	}
	if err := db.Delete(&oBids, delBids).Error; err != nil {
		return err
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
