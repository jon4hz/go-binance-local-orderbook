package binance_futures

import (
	"context"
	"fmt"
	"log"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/jon4hz/go-binance-local-orderbook/config"
	"github.com/jon4hz/go-binance-local-orderbook/database"
	"github.com/jon4hz/go-binance-local-orderbook/exchange"
)

func HandleWebsocket(config *config.Config) {
	var response database.DatabaseInsert
	wsDepthHandler := func(event *futures.WsDepthEvent) {
		response = &database.BinanceFuturesDepthResponse{Response: event}
		exchange.BigU = event.FirstUpdateID
		exchange.SmallU = event.LastUpdateID
		// first time
		if exchange.Prev_u == 0 {
			exchange.Prev_u = exchange.SmallU
			snap, err := downloadSnapshot(*config)
			response = &database.BinanceFuturesDepthResponse{Snapshot: snap}
			if err != nil {
				panic("Error while downloading the snapshot")
			}
			response.InsertIntoDatabase()
			exchange.LastUpdateID = snap.LastUpdateID
			fmt.Println(exchange.LastUpdateID)

		} else {
			response.InsertIntoDatabase()
			fmt.Println(exchange.SmallU, exchange.Prev_u+1, exchange.BigU)
			exchange.Prev_u = exchange.SmallU
		}

	}
	errHandler := func(err error) {
		log.Fatal(err)
	}
	doneC, _, err := futures.WsDiffDepthServe(config.Exchange.Market, wsDepthHandler, errHandler)
	if err != nil {
		log.Fatal(err)
		return
	}
	<-doneC
}

func downloadSnapshot(config config.Config) (res *futures.DepthResponse, err error) {
	client := binance.NewFuturesClient("", "")
	res, err = client.NewDepthService().Symbol(config.Exchange.Market).
		Limit(1000).
		Do(context.Background())
	return

}
