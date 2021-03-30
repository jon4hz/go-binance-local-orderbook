// this package is heavily inspired by github.com/TwinProduction/gatus

package config

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"

	"github.com/jon4hz/go-binance-local-orderbook/database"
	"github.com/jon4hz/go-binance-local-orderbook/exchange"
	"github.com/spf13/viper"
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

	config *Config
)

type Config struct {
	Exchange *exchange.Config `mapstructure:"exchange"`
	Database *database.Config `mapstructure:"database"`
}

func Get() *Config {
	if config == nil {
		panic(ErrConfigNotLoaded)
	}
	return config
}

func Load(configFile string) error {
	cfg, err := readConfiguration(configFile)
	if err != nil {
		return err
	}
	config = cfg
	return nil
}

func readConfiguration(fileName string) (config *Config, err error) {

	viper.SetConfigType("yaml")

	// check if file exists
	var readFromFile bool

	var succ fs.FileInfo
	succ, err = os.Stat(fileName)
	if succ != nil && !(os.IsNotExist(err)) {
		viper.SetConfigFile(fileName)
		log.Printf("[config][Load] Reading configuration from configFile=%s", fileName)
		readFromFile = true
	} else {
		readFromFile = false
		log.Print("[config][Load] Reading configuration from environment vars")
	}

	// set defaults
	viper.SetDefault("database.POSTGRES_PORT", "5432")

	// map environment variables to yaml values
	viper.BindEnv("exchange.NAME", "NAME")
	viper.BindEnv("exchange.MARKET", "MARKET")

	viper.BindEnv("database.POSTGRES_DB", "POSTGRES_DB")
	viper.BindEnv("database.POSTGRES_USERNAME", "POSTGRES_USERNAME")
	viper.BindEnv("database.POSTGRES_PASSWORD", "POSTGRES_PASSWORD")
	viper.BindEnv("database.POSTGRES_SERVER", "POSTGRES_SERVER")
	viper.BindEnv("database.POSTGRES_PORT", "POSTGRES_PORT")

	viper.AutomaticEnv()

	if readFromFile {
		err = viper.ReadInConfig()
		if err != nil {
			return
		}
	}

	err = viper.Unmarshal(&config)

	if err == nil {
		validateExchangeConfig(config)
		validateDatabaseConfig(config)
	}

	return
}

func validateExchangeConfig(config *Config) {
	if config.Exchange == nil {
		panic("[config][validateExchangeConfig] Exchange is not configured")
	}
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
	if config.Database.DBUsername == "" {
		panic("[config][validateDatabaseConfig] Database Username is not configured")
	}
	if config.Database.DBServer == "" {
		panic("[config][validateDatabaseConfig] Database Server is not configured")
	}
}
