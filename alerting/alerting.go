package alerting

import (
	"log"

	"github.com/jon4hz/go-binance-local-orderbook/alerting/telegram"
)

type AlertingMSG string

type Config struct {
	Telegram *telegram.Config `mapstructure:"telegram"`
}

func (msg AlertingMSG) TriggerAlert(cfg *Config) {
	if cfg == nil {
		log.Printf("[alerting][Trigger] No alerting provider configured")
		return
	}
	if cfg.Telegram != nil {
		if err := telegram.TriggerTelegramAlert(cfg.Telegram, msg); err != nil {
			log.Printf("[alerting][Telegram] Error: %s\n", err)
		}
	}
}
