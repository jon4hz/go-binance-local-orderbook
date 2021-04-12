package binance_futures

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/jon4hz/go-binance-local-orderbook/config"
	"github.com/jon4hz/go-binance-local-orderbook/database"
	"github.com/jon4hz/go-binance-local-orderbook/exchange"
)

func HandleWebsocket(config *config.Config) {
	wsDepthHandler := func(event *futures.WsDepthEvent) {
		response := &database.BinanceFuturesDepthResponse{Response: event}
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
			err = response.InsertIntoDatabase(config.Database.DBTableMarketName)
			if err != nil {
				log.Println(err)
				// send notification
			}
			exchange.LastUpdateID = snap.LastUpdateID
			fmt.Println(exchange.LastUpdateID)

		} else {
			err := response.InsertIntoDatabase(config.Database.DBTableMarketName)
			if err != nil {
				log.Println(err)
				// send notification
			}
			fmt.Println(exchange.SmallU, exchange.Prev_u+1, exchange.BigU)
			exchange.Prev_u = exchange.SmallU
		}

	}
	errHandler := func(err error) {
		log.Fatal(err)
	}
	var monitorWS func(sym string, ch chan struct{})
	monitorWS = func(sym string, ch chan struct{}) {
		go func() {
			<-ch
			// ws disconnected, try to re-establish.
			log.Printf("Websocket for %s crashed, spawning a new one.", sym)
			doneC, _, err := futures.WsDiffDepthServe(sym, wsDepthHandler, errHandler)
			if err != nil {
				log.Printf("error registering symbol %s: %v", sym, err)
				return
			}
			monitorWS(sym, doneC)

			<-doneC
		}()
	}
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func(sym string) {
		defer wg.Done()
		doneC, _, err := futures.WsDiffDepthServe(sym, wsDepthHandler, errHandler)
		if err != nil {
			log.Printf("error registering symbol %s: %v", sym, err)
		}
		monitorWS(sym, doneC)

		<-doneC
	}(config.Exchange.Market)
	wg.Wait()
}

func downloadSnapshot(config config.Config) (res *futures.DepthResponse, err error) {
	client := binance.NewFuturesClient("", "")
	res, err = client.NewDepthService().Symbol(config.Exchange.Market).
		Limit(1000).
		Do(context.Background())
	return

}
