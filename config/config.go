// this package is heavily inspired by github.com/TwinProduction/gatus

package config

import (
	"errors"
	"log"
	"os"

	"github.com/jon4hz/go-binance-local-orderbook/database"
	"github.com/jon4hz/go-binance-local-orderbook/exchange"
	"github.com/spf13/viper"
)

const (
	DefaultConfigurationFilePath = "config/config.yml"

	DefaultFallbackConfigurationFilePath = "config/config.yml"
)

var (
	// ErrNoServiceInConfig is an error returned when a configuration file has no services configured
	ErrNoServiceInConfig = errors.New("configuration file should contain at least 1 service")

	// ErrConfigFileNotFound is an error returned when the configuration file could not be found
	ErrConfigFileNotFound = errors.New("configuration file not found")

	// ErrConfigNotLoaded is an error returned when an attempt to Get() the configuration before loading it is made
	ErrConfigNotLoaded = errors.New("configuration is nil")

	// ErrInvalidSecurityConfig is an error returned when the security configuration is invalid
	ErrInvalidSecurityConfig = errors.New("invalid security configuration")
)

type Config struct {
	Exchange *exchange.Config `mapstructure:"exchange"`
	Database *database.Config `mapstructure:"database"`
}

func Load(configFile string) (Config, error) {
	log.Printf("[config][Load] Reading configuration from configFile=%s", configFile)
	cfg, err := readConfiguration(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, ErrConfigFileNotFound
		}
		return cfg, err
	}
	return cfg, nil
}

func readConfiguration(fileName string) (config Config, err error) {
	viper.SetConfigFile(fileName)
	viper.SetConfigType("yaml")

	// set defaults
	viper.SetDefault("database.POSTGRES_PORT", "5432")

	// map environment variables to yaml values
	viper.BindEnv("exchange.NAME", "NAME")
	viper.BindEnv("exchange.MARKET", "MARKET")
	viper.BindEnv("exchange.API_KEY", "API_KEY")
	viper.BindEnv("exchange.API_SECRET", "API_SECRET")

	viper.BindEnv("database.POSTGRES_DATABASE", "POSTGRES_DATABASE")
	viper.BindEnv("database.POSTGRES_USERNAME", "POSTGRES_USERNAME")
	viper.BindEnv("database.POSTGRES_PASSWORD", "POSTGRES_PASSWORD")
	viper.BindEnv("database.POSTGRES_SERVER", "POSTGRES_SERVER")
	viper.BindEnv("database.POSTGRES_PORT", "POSTGRES_PORT")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)

	validateExchangeConfig(config)
	validateDatabaseConfig(config)

	return
}

func validateExchangeConfig(config Config) {
	if config.Exchange == nil {
		panic("[config][validateExchangeConfig] Exchange is not configured")
	}
	if config.Exchange.APIKey == "" {
		panic("[config][validateExchangeConfig] Exchange API Key is not configured")
	}
	if config.Exchange.APISecret == "" {
		panic("[config][validateExchangeConfig] Exchange API Key is not configured")
	}
	if config.Exchange.Name == "" {
		panic("[config][validateExchangeConfig] Exchange Name is not configured")
	}
	if config.Exchange.Market == "" {
		panic("[config][validateExchangeConfig] Exchange Market is not configured")
	}
}

func validateDatabaseConfig(config Config) {
	if config.Database == nil {
		panic("[config][validateDatabaseConfig] Database is not configured")
	}
	if config.Database.DBName == "" {
		panic("[config][validateDatabaseConfig] Database Name is not configured")
	}
	if config.Database.DBPassword == "" {
		panic("[config][validateDatabaseConfig] Database Password is not configured")
	}
	if config.Database.DBPort == "" {
		panic("[config][validateDatabaseConfig] Database Port is not configured")
	}
	if config.Database.DBUsername == "" {
		panic("[config][validateDatabaseConfig] Database Username is not configured")
	}
	if config.Database.DBServer == "" {
		panic("[config][validateDatabaseConfig] Database Server is not configured")
	}
}