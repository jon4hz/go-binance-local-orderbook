package config

type ExchangeConfig struct {
	Name   string `mapstructure:"NAME"`
	Market string `mapstructure:"MARKET"`
}
