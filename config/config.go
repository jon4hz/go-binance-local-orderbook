package config

import (
	"errors"
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jon4hz/go-binance-local-orderbook/database"
	"github.com/jon4hz/go-binance-local-orderbook/exchange"
)

const (
	DefaultConfigurationFilePath = "config/config.yml"
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
	Exchange      *exchange.Config `yaml:"exchange"`
	Database      *database.Config `yaml:"database"`
	DeleteOldSnap bool             `yaml:"deleteOldSnap" env:"DeleteOldSnap" env-default:"true"`
}

func Load(configFile string) (cfg Config, err error) {
	err = cleanenv.ReadConfig(configFile, &cfg)
	if err != nil {
		return
	}
	if err == nil {
		validateExchangeConfig(&cfg)
		validateDatabaseConfig(&cfg)
		validateOtherConfig(&cfg)
	}
	return

}

func validateExchangeConfig(config *Config) {
	/* if config.Exchange == nil {
		panic("[config][validateExchangeConfig] Exchange is not configured")
	} */
	if config.Exchange.Name == "" {
		panic("[config][validateExchangeConfig] Exchange Name is not configured")
	} else {
		switch config.Exchange.Name {
		case
			"binance",
			"binance-futures":
			// pass
		default:
			panic(fmt.Sprintf("[config][validateExchangeConfig] Exchange Name can't be %s", config.Exchange.Name))
		}
	}
	if config.Exchange.Market == "" {
		panic("[config][validateExchangeConfig] Exchange Market is not configured")
	}
}

func validateDatabaseConfig(config *Config) {
	// config.Database always exists, since config.Database.Port has a default value
	/* if config.Database == nil {
		panic("[config][validateDatabaseConfig] Database is not configured")
	} */
	if config.Database.DBName == "" {
		panic("[config][validateDatabaseConfig] Database Name is not configured")
	}
	if config.Database.DBPassword == "" {
		panic("[config][validateDatabaseConfig] Database Password is not configured")
	}
	// config.Database.DBPort has a default value and can't be ""
	/* if config.Database.DBPort == "" {
		panic("[config][validateDatabaseConfig] Database Port is not configured")
	} */
	if config.Database.DBUser == "" {
		panic("[config][validateDatabaseConfig] Database User is not configured")
	}
	if config.Database.DBServer == "" {
		panic("[config][validateDatabaseConfig] Database Server is not configured")
	}
}

func validateOtherConfig(config *Config) {
	// Will never happen, DeleteOldSnap has default value (true)
	/* if !config.DeleteOldSnap {
		panic("[config][validateOtherConfig] DeleteOldSnap is not configured")
	} */
}
