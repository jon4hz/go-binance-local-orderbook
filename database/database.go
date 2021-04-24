package database

import (
	"fmt"
	"log"
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
	log.Println("[database][migrator] GORM migration successful")
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

func doDBInsert(sym string, asks []ask, bids []bid) error {

	if err := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&asks).Error; err != nil {
		return err
	}

	if err := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&bids).Error; err != nil {
		return err
	}

	// loop again over bids and asks to delete 0 values
	var delAsks = []float64{}
	for _, v := range asks {
		if v.Quantity == 0 {
			delAsks = append(delAsks, v.Price)
		}
	}
	if len(delAsks) > 0 {
		if err := db.Delete(&asks, delAsks).Error; err != nil {
			return err
		}
	}

	var delBids = []float64{}
	for _, v := range bids {
		if v.Quantity == 0 {
			delBids = append(delBids, v.Price)
		}
	}
	if len(delBids) > 0 {
		if err := db.Delete(&bids, delBids).Error; err != nil {
			return err
		}
	}

	return nil
}

func (resp *BinanceDepthResponse) InsertIntoDatabase(sym string) error {
	var err error
	if resp.Snapshot != nil {
		// convert asks to int64
		asks := make([]ask, len(resp.Snapshot.Asks))
		for i, v := range resp.Snapshot.Asks {
			asks[i].Price, err = strconv.ParseFloat(v.Price, 64)
			if err != nil {
				return fmt.Errorf("error converting ask price string to int64: %s", err)
			}
			asks[i].Quantity, err = strconv.ParseFloat(v.Quantity, 64)
			if err != nil {
				return fmt.Errorf("error converting ask quantity string to int64: %s", err)
			}
		}
		// convert bids to int64
		bids := make([]bid, len(resp.Snapshot.Bids))
		for i, v := range resp.Snapshot.Bids {
			bids[i].Price, err = strconv.ParseFloat(v.Price, 64)
			if err != nil {
				return fmt.Errorf("error converting bid price string to int64: %s", err)
			}
			bids[i].Quantity, err = strconv.ParseFloat(v.Quantity, 64)
			if err != nil {
				return fmt.Errorf("error converting bid quantity string to int64: %s", err)
			}
		}
		// insert into db
		err = doDBInsert(sym, asks, bids)
		if err != nil {
			return err
		}
	}

	if resp.Response != nil {
		// convert asks to int64
		asks := make([]ask, len(resp.Response.Asks))
		for i, v := range resp.Response.Asks {
			asks[i].Price, err = strconv.ParseFloat(v.Price, 64)
			if err != nil {
				return fmt.Errorf("error converting ask price string to int64: %s", err)
			}
			asks[i].Quantity, err = strconv.ParseFloat(v.Quantity, 64)
			if err != nil {
				return fmt.Errorf("error converting ask quantity string to int64: %s", err)
			}
		}
		// convert bids to int64
		bids := make([]bid, len(resp.Response.Bids))
		for i, v := range resp.Response.Bids {
			bids[i].Price, err = strconv.ParseFloat(v.Price, 64)
			if err != nil {
				return fmt.Errorf("error converting bid price string to int64: %s", err)
			}
			bids[i].Quantity, err = strconv.ParseFloat(v.Quantity, 64)
			if err != nil {
				return fmt.Errorf("error converting bid quantity string to int64: %s", err)
			}
		}
		// insert into db
		err = doDBInsert(sym, asks, bids)
		if err != nil {
			return err
		}
	}
	return nil
}

func (resp *BinanceFuturesDepthResponse) InsertIntoDatabase(sym string) error {
	var err error
	if resp.Snapshot != nil {
		// convert asks to int64
		asks := make([]ask, len(resp.Snapshot.Asks))
		for i, v := range resp.Snapshot.Asks {
			asks[i].Price, err = strconv.ParseFloat(v.Price, 64)
			if err != nil {
				return fmt.Errorf("error converting ask price string to int64: %s", err)
			}
			asks[i].Quantity, err = strconv.ParseFloat(v.Quantity, 64)
			if err != nil {
				return fmt.Errorf("error converting ask quantity string to int64: %s", err)
			}
		}
		// convert bids to int64
		bids := make([]bid, len(resp.Snapshot.Bids))
		for i, v := range resp.Snapshot.Bids {
			bids[i].Price, err = strconv.ParseFloat(v.Price, 64)
			if err != nil {
				return fmt.Errorf("error converting bid price string to int64: %s", err)
			}
			bids[i].Quantity, err = strconv.ParseFloat(v.Quantity, 64)
			if err != nil {
				return fmt.Errorf("error converting bid quantity string to int64: %s", err)
			}
		}
		// insert into db
		err = doDBInsert(sym, asks, bids)
		if err != nil {
			return err
		}
	}

	if resp.Response != nil {
		// convert asks to int64
		asks := make([]ask, len(resp.Response.Asks))
		for i, v := range resp.Response.Asks {
			asks[i].Price, err = strconv.ParseFloat(v.Price, 64)
			if err != nil {
				return fmt.Errorf("error converting ask price string to int64: %s", err)
			}
			asks[i].Quantity, err = strconv.ParseFloat(v.Quantity, 64)
			if err != nil {
				return fmt.Errorf("error converting ask quantity string to int64: %s", err)
			}
		}
		// convert bids to int64
		bids := make([]bid, len(resp.Response.Bids))
		for i, v := range resp.Response.Bids {
			bids[i].Price, err = strconv.ParseFloat(v.Price, 64)
			if err != nil {
				return fmt.Errorf("error converting bid price string to int64: %s", err)
			}
			bids[i].Quantity, err = strconv.ParseFloat(v.Quantity, 64)
			if err != nil {
				return fmt.Errorf("error converting bid quantity string to int64: %s", err)
			}
		}
		// insert into db
		err = doDBInsert(sym, asks, bids)
		if err != nil {
			return err
		}
	}
	return nil
}
