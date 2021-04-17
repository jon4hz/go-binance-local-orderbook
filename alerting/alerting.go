package alerting

import (
	"errors"

	"github.com/jon4hz/go-binance-local-orderbook/alerting/telegram"
	"github.com/jon4hz/go-binance-local-orderbook/config"
)

type AlertingMSG string

func TriggerAlert(cfg *config.Config, msg AlertingMSG) error {
	if cfg.Alerting == nil {
		return errors.New("[alerting][Trigger] No alerting provider configured")
	}
	if cfg.Alerting.Telegram != nil {
		if err := telegram.TriggerTelegramAlert(cfg, msg); err != nil {
			return err
		}
	}
	return nil
}
