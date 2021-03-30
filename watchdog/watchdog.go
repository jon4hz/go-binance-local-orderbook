package watchdog

import (
	"fmt"
	"time"

	"github.com/jon4hz/go-binance-local-orderbook/handler"
)

func Watcher() {
	var prev_u int64
	for {
		time.Sleep(50 * time.Second)
		if prev_u == handler.BigU {
			fmt.Println("Error: orderbook didn't change for 50 seconds.")
		}
		prev_u = handler.BigU
	}
}
