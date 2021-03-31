package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Config struct {
	DBName     string `yaml:"db" env:"POSTGRES_DB"`
	DBUser     string `yaml:"user" env:"POSTGRES_USER"`
	DBPassword string `yaml:"password" env:"POSTGRES_PASSWORD"`
	DBServer   string `yaml:"server" env:"POSTGRES_SERVER"`
	DBPort     string `yaml:"port" env:"POSTGRES_PORT" env-default:"5432"`
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
