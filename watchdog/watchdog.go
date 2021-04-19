package watchdog

import (
	"log"
	"time"

	"github.com/jon4hz/go-binance-local-orderbook/alerting"
	"github.com/jon4hz/go-binance-local-orderbook/config"
	"github.com/jon4hz/go-binance-local-orderbook/exchange"
)

func Watcher(cfg *config.Config) {
	var prev_u int64
	for {
		time.Sleep(50 * time.Second)
		if prev_u == exchange.SmallU {
			log.Println("Error: orderbook didn't change for 50 seconds.")
			msg := alerting.AlertingMSG("🚨 Error: orderbook didn't change for 50 seconds.")
			msg.TriggerAlert(cfg)
		}
		prev_u = exchange.SmallU
	}
}
