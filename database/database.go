package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Config struct {
	DBName     string `mapstructure:"POSTGRES_DB"`
	DBUser     string `mapstructure:"POSTGRES_USER"`
	DBPassword string `mapstructure:"POSTGRES_PASSWORD"`
	DBServer   string `mapstructure:"POSTGRES_SERVER"`
	DBPort     string `mapstructure:"POSTGRES_PORT"`
}

type testRows struct {
	name string
}

func InitDatabase(dbpool *pgxpool.Pool) (err error) {

	rows, err := dbpool.Query(context.Background(), "SELECT tablename FROM pg_catalog.pg_tables;")
	if err != nil {
		log.Fatal("Unable to connect to database: ", err)
		os.Exit(1)
	}
	defer rows.Close()
	var test string
	rows.Scan(&test)
	fmt.Println(test)
	return
}
