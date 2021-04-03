package watchdog

import (
	"fmt"
	"time"

	"github.com/jon4hz/go-binance-local-orderbook/exchange"
)

func Watcher() {
	var prev_u int64
	for {
		time.Sleep(50 * time.Second)
		if prev_u == exchange.SmallU {
			fmt.Println("Error: orderbook didn't change for 50 seconds.")
		}
		prev_u = exchange.SmallU
	}
}
