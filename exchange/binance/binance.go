package binance

import (
	"context"
	"fmt"

	"github.com/adshao/go-binance/v2"
	"github.com/jon4hz/go-binance-local-orderbook/config"
	"github.com/jon4hz/go-binance-local-orderbook/exchange"
)

func HandleWebsocket(config *config.Config) {
	wsDepthHandler := func(event *binance.WsDepthEvent) {
		exchange.BigU = event.FirstUpdateID
		exchange.SmallU = event.UpdateID
		// first time
		if exchange.Prev_u == 0 {
			exchange.Prev_u = exchange.SmallU
			snap, err := downloadSnapshot(*config)
			if err != nil {
				panic("Error while downloading the snapshot")
			}
			exchange.LastUpdateID = snap.LastUpdateID
			fmt.Println(exchange.LastUpdateID)

		}
		fmt.Println(exchange.SmallU, exchange.Prev_u+1, exchange.BigU)
		exchange.Prev_u = exchange.SmallU

	}
	errHandler := func(err error) {
		fmt.Println(err)
	}
	doneC, _, err := binance.WsDepthServe(config.Exchange.Market, wsDepthHandler, errHandler)
	if err != nil {
		fmt.Println(err)
		return
	}
	<-doneC
}

func downloadSnapshot(config config.Config) (res *binance.DepthResponse, err error) {
	client := binance.NewClient("", "")
	res, err = client.NewDepthService().Symbol(config.Exchange.Market).
		Do(context.Background())
	return

}
