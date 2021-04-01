package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Config struct {
	DBName            string `mapstructure:"POSTGRES_DB"`
	DBUser            string `mapstructure:"POSTGRES_USER"`
	DBPassword        string `mapstructure:"POSTGRES_PASSWORD"`
	DBServer          string `mapstructure:"POSTGRES_SERVER"`
	DBPort            string `mapstructure:"POSTGRES_PORT"`
	DBTableMarketName string
	DBDeleteOldSnap   bool
}

func InitDatabase(dbpool *pgxpool.Pool, ctx context.Context, config *Config) error {
	// drop old tables if set to true (default)
	conn, err := dbpool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()
	if config.DBDeleteOldSnap {
		tables := [3]string{"asks", "bids", "general"}
		for _, table := range tables {
			if _, err := conn.Exec(ctx,
				fmt.Sprintf(`DROP TABLE IF EXISTS %s_%s;`, config.DBTableMarketName, table)); err != nil {
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
		);`, config.DBTableMarketName, config.DBTableMarketName, config.DBTableMarketName)); err != nil {
		return err
	}
	log.Println("Successfully created new tables")
	return nil
}
