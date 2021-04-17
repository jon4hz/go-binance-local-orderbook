package alerting

import (
	"log"

	"github.com/jon4hz/go-binance-local-orderbook/alerting/telegram"
	"github.com/jon4hz/go-binance-local-orderbook/config"
)

type AlertingMSG string

func (msg AlertingMSG) TriggerAlert(cfg *config.Config) {
	if cfg.Alerting == nil {
		log.Printf("[alerting][Trigger] No alerting provider configured")
		return
	}
	if cfg.Alerting.Telegram != nil {
		if err := telegram.TriggerTelegramAlert(cfg, msg); err != nil {
			log.Printf("[alerting][Telegram] Error: %s\n", err)
		}
	}
}
