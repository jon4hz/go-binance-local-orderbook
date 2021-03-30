package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
)

type Config struct {
	DBName     string `mapstructure:"POSTGRES_DB"`
	DBUsername string `mapstructure:"POSTGRES_USER"`
	DBPassword string `mapstructure:"POSTGRES_PASSWORD"`
	DBServer   string `mapstructure:"POSTGRES_SERVER"`
	DBPort     string `mapstructure:"POSTGRES_PORT"`
}

func InitDatabase(*Config) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
}
