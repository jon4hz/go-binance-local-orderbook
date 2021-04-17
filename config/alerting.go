package config

type AlertingConfig struct {
	Telegram *TelegramConfig `mapstructure:"telegram"`
}

type TelegramConfig struct {
	Token string `mapstructure:"TOKEN"`
	Chat  int64  `mapstructure:"CHAT"`
}
