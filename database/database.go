package database

import (
	"fmt"
	"log"
	"reflect"
	"strconv"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

type Config struct {
	DBName            string `mapstructure:"POSTGRES_DB"`
	DBUser            string `mapstructure:"POSTGRES_USER"`
	DBPassword        string `mapstructure:"POSTGRES_PASSWORD"`
	DBServer          string `mapstructure:"POSTGRES_SERVER"`
	DBPort            string `mapstructure:"POSTGRES_PORT"`
	Debug             bool   `mapstructure:"Debug"`
	DBTableMarketName string
	DBDeleteOldSnap   bool
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

	var gormConfig = &gorm.Config{}

	if !config.Debug {
		gormConfig.Logger = logger.Default.LogMode(logger.Silent)
		log.Println("[database][logger] GORM Logger is disabled")
	} else {
		log.Println("[database][logger] GORM Logger is enabled")
	}

	var err error
	db, err = gorm.Open(postgres.Open(dsn), gormConfig)
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
		log.Println("[database][migrator] Deleted old tables successfully")
	}
	err = db.AutoMigrate(&ask{}, &bid{})
	log.Println("[database][migrator] GORM migration successfull")
	return
}

type bid struct {
	Price    float64 `gorm:"primaryKey"`
	Quantity float64
}

type ask struct {
	Price    float64 `gorm:"primaryKey"`
	Quantity float64
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

func createUnifiedStruct(asks, bids interface{}) ([]ask, []bid, error) {
	var err error

	tmpA := reflect.ValueOf(asks)
	// convert datatype
	if tmpA.Kind() != reflect.Slice {
		return nil, nil, fmt.Errorf("error converting asks interface to slice")
	}

	oAsks := make([]ask, tmpA.Len())

	for i := 0; i < tmpA.Len(); i++ {
		switch x := tmpA.Index(i).Interface().(type) {
		case binance.Ask:
			oAsks[i].Price, err = strconv.ParseFloat(x.Price, 64)
			if err != nil {
				return nil, nil, fmt.Errorf("error converting ask string to int64")
			}
		case futures.Ask:
			oAsks[i].Price, err = strconv.ParseFloat(x.Price, 64)
			if err != nil {
				return nil, nil, fmt.Errorf("error converting ask string to int64")
			}
		default:
			return nil, nil, fmt.Errorf("error no matching type found for switch statement")
		}

	}

	tmpB := reflect.ValueOf(bids)
	// convert datatype
	if tmpB.Kind() != reflect.Slice {
		return nil, nil, fmt.Errorf("error converting bids interface to slice")
	}

	oBids := make([]bid, tmpB.Len())

	for i := 0; i < tmpB.Len(); i++ {
		switch x := tmpB.Index(i).Interface().(type) {
		case binance.Bid:
			oBids[i].Price, err = strconv.ParseFloat(x.Price, 64)
			if err != nil {
				return nil, nil, fmt.Errorf("error converting bid string to int64")
			}
		case futures.Bid:
			oBids[i].Price, err = strconv.ParseFloat(x.Price, 64)
			if err != nil {
				return nil, nil, fmt.Errorf("error converting bid string to int64")
			}
		default:
			return nil, nil, fmt.Errorf("error no matching type found for switch statement")
		}

	}

	return oAsks, oBids, nil
}

func doDBInsert(sym string, asks interface{}, bids interface{}) error {
	oAsks, oBids, err := createUnifiedStruct(asks, bids)
	if err != nil {
		return err
	}

	if err := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&oAsks).Error; err != nil {
		return err
	}

	if err := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&oBids).Error; err != nil {
		return err
	}

	// loop again over bids and asks to delete 0 values
	var delAsks = []string{}
	/* for _, v := range oAsks {
		var quant float64
		if quant, err = strconv.ParseFloat(v.Quantity, 64); err != nil {
			log.Printf("[database][dbinsert] couldn't convert \"quantity\" to float: %s\n", err)
			return err
		}
		if quant == 0 {
			delAsks = append(delAsks, v.Price)
		}
	} */
	if err := db.Delete(&oAsks, delAsks).Error; err != nil {
		return err
	}

	var delBids = []string{}
	/* for _, v := range oBids {
		var quant float64
		if quant, err = strconv.ParseFloat(v.Quantity, 64); err != nil {
			log.Printf("[database][dbinsert] couldn't convert \"quantity\" to float: %s\n", err)
			return err
		}
		if quant == 0 {
			delBids = append(delBids, v.Price)
		}
	} */
	if err := db.Delete(&oBids, delBids).Error; err != nil {
		return err
	}

	return nil
}

func (resp *BinanceDepthResponse) InsertIntoDatabase(sym string) error {
	if resp.Snapshot != nil {
		err := doDBInsert(sym, resp.Snapshot.Asks, resp.Snapshot.Bids)
		if err != nil {
			return err
		}
	}
	if resp.Response != nil {
		err := doDBInsert(sym, resp.Response.Asks, resp.Response.Bids)
		if err != nil {
			return err
		}
	}
	return nil
}

func (resp *BinanceFuturesDepthResponse) InsertIntoDatabase(sym string) error {
	if resp.Snapshot != nil {
		err := doDBInsert(sym, resp.Snapshot.Asks, resp.Snapshot.Bids)
		if err != nil {
			return err
		}
	}
	if resp.Response != nil {
		err := doDBInsert(sym, resp.Response.Asks, resp.Response.Bids)
		if err != nil {
			return err
		}
	}
	return nil
}
