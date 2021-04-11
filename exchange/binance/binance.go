package binance

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/adshao/go-binance/v2"
	"github.com/jon4hz/go-binance-local-orderbook/config"
	"github.com/jon4hz/go-binance-local-orderbook/exchange"
)

func InitWebsocket(config *config.Config) {
	//var response database.DatabaseInsert
	wsDepthHandler := func(event *binance.WsDepthEvent) {
		//response = &database.BinanceDepthResponse{Response: event}
		exchange.BigU = event.FirstUpdateID
		exchange.SmallU = event.UpdateID
		// first time
		if exchange.Prev_u == 0 {
			exchange.Prev_u = exchange.SmallU
			snap, err := downloadSnapshot(*config)
			//response = &database.BinanceDepthResponse{Snapshot: snap}
			if err != nil {
				panic("Error while downloading the snapshot")
			}
			//response.InsertIntoDatabase()
			exchange.LastUpdateID = snap.LastUpdateID
			fmt.Println(exchange.LastUpdateID)

		} else {
			//response.InsertIntoDatabase()
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
			doneC, _, err := binance.WsDepthServe(sym, wsDepthHandler, errHandler)
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
		doneC, _, err := binance.WsDepthServe(sym, wsDepthHandler, errHandler)
		if err != nil {
			log.Printf("error registering symbol %s: %v", sym, err)
		}
		monitorWS(sym, doneC)

		<-doneC
	}(config.Exchange.Market)
	wg.Wait()

}

func downloadSnapshot(config config.Config) (res *binance.DepthResponse, err error) {
	client := binance.NewClient("", "")
	res, err = client.NewDepthService().Symbol(config.Exchange.Market).
		Limit(1000).
		Do(context.Background())
	return

}
