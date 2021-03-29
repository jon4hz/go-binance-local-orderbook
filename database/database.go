package database

type Config struct {
	DBName     string `mapstructure:"POSTGRES_DATABASE"`
	DBUsername string `mapstructure:"POSTGRES_USERNAME"`
	DBPassword string `mapstructure:"POSTGRES_PASSWORD"`
	DBServer   string `mapstructure:"POSTGRES_SERVER"`
	DBPort     string `mapstructure:"POSTGRES_PORT"`
}
