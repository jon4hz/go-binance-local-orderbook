package config

type DatabaseConfig struct {
	DBName            string `mapstructure:"POSTGRES_DB"`
	DBUser            string `mapstructure:"POSTGRES_USER"`
	DBPassword        string `mapstructure:"POSTGRES_PASSWORD"`
	DBServer          string `mapstructure:"POSTGRES_SERVER"`
	DBPort            string `mapstructure:"POSTGRES_PORT"`
	DBTableMarketName string
	DBDeleteOldSnap   bool
}
