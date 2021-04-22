package binance_futures

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/jon4hz/go-binance-local-orderbook/alerting"
	"github.com/jon4hz/go-binance-local-orderbook/config"
	"github.com/jon4hz/go-binance-local-orderbook/database"
	"github.com/jon4hz/go-binance-local-orderbook/exchange"
)

func InitWebsocket(config *config.Config) {

	// set keepalive vars
	binance.WebsocketKeepalive = true
	binance.WebsocketTimeout = time.Second * 6

	//var response database.DatabaseInsert
	wsDepthHandler := func(event *futures.WsDepthEvent) {
		exchange.Notified = false
		response := &database.BinanceFuturesDepthResponse{Response: event}
		exchange.BigU = event.FirstUpdateID
		exchange.SmallU = event.LastUpdateID

		// first time
		if exchange.Prev_u == 0 {
			// download snapshot
			if exchange.LastUpdateID == 0 {
				snap, err := downloadSnapshot(config)
				if err != nil {
					log.Println("[exchange][snapshot] Error while downloading the snapshot: ", err)
					return
				}
				response = &database.BinanceFuturesDepthResponse{Snapshot: snap}
				err = response.InsertIntoDatabase(config.Database.DBTableMarketName)
				if err != nil {
					log.Println(err)
					msg := alerting.AlertingMSG(fmt.Sprintf("🚨 Error coin: %s, couldn't insert data %s", config.Exchange.Market, err))
					go msg.TriggerAlert(config.Alerting)
					return
				}
				msg := alerting.AlertingMSG(fmt.Sprintf("💡 Info: Downloaded new snapshot for coin: %s", config.Exchange.Market))
				go msg.TriggerAlert(config.Alerting)
				log.Println("[exchange][dbinsert] Inserted snapshot into db")
				exchange.LastUpdateID = snap.LastUpdateID
			}
			if exchange.SmallU >= exchange.LastUpdateID+1 && exchange.BigU <= exchange.LastUpdateID+1 {
				err := response.InsertIntoDatabase(config.Database.DBTableMarketName)
				if err != nil {
					log.Println(err)
					msg := alerting.AlertingMSG(fmt.Sprintf("🚨 Error coin: %s, couldn't insert data %s", config.Exchange.Market, err))
					go msg.TriggerAlert(config.Alerting)
					return
				}
				exchange.Prev_u = exchange.SmallU
				log.Println("[exchange][dbinsert] Inserted first event successfully")
			}
			return

		} else if event.PrevLastUpdateID >= exchange.Prev_u {
			if event.PrevLastUpdateID > exchange.Prev_u {
				log.Printf("[exchange][missmatch cond.] Warning, U = %d and prev_u = %d", event.PrevLastUpdateID, exchange.Prev_u)
				msg := alerting.AlertingMSG(fmt.Sprintf("⚠️ Warning: Orderbook could be out of sync for %s, U: %d, prev_u: %d", config.Exchange.Market, event.PrevLastUpdateID, exchange.Prev_u))
				go msg.TriggerAlert(config.Alerting)
			}
			err := response.InsertIntoDatabase(config.Database.DBTableMarketName)
			if err != nil {
				log.Println(err)

				msg := alerting.AlertingMSG(fmt.Sprintf("🚨 Error coin: %s, couldn't insert data %s", config.Exchange.Market, err))
				go msg.TriggerAlert(config.Alerting)
			}
			exchange.Prev_u = exchange.SmallU
		} else {
			log.Println("Error")
		}
	}
	errHandler := func(err error) {
		log.Printf("Error: %s", err)
		msg := alerting.AlertingMSG(fmt.Sprintf("🚨 Websocket error: %s", err))
		go msg.TriggerAlert(config.Alerting)
	}
	var monitorWS func(sym string, ch chan struct{})
	monitorWS = func(sym string, ch chan struct{}) {
		go func() {
			<-ch
			// ws disconnected, try to re-establish.
			log.Printf("[exchange][websocket] Websocket for %s crashed, spawning a new one.", sym)
			doneC, _, err := futures.WsDiffDepthServe(sym, wsDepthHandler, errHandler)
			if err != nil {
				log.Printf("[exchange][reconnect] Error registering symbol %s: %v", sym, err)
				for err != nil {
					time.Sleep(time.Second / 10)
					doneC, _, err = futures.WsDiffDepthServe(sym, wsDepthHandler, errHandler)

				}
				log.Println("[exchange][reconnect] Connection established", sym, err)
			}
			monitorWS(sym, doneC)
		}()
	}
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func(sym string) {
		defer wg.Done()
		doneC, _, err := futures.WsDiffDepthServe(sym, wsDepthHandler, errHandler)
		if err != nil {
			log.Printf("[exchange][reconnect] Error registering symbol %s: %v", sym, err)
			for err != nil {
				time.Sleep(time.Second / 10)
				doneC, _, err = futures.WsDiffDepthServe(sym, wsDepthHandler, errHandler)
			}
		}
		monitorWS(sym, doneC)

	}(config.Exchange.Market)
	wg.Wait()
}

func downloadSnapshot(config *config.Config) (res *futures.DepthResponse, err error) {
	client := futures.NewClient("", "")
	res, err = client.NewDepthService().Symbol(config.Exchange.Market).
		Limit(1000).
		Do(context.TODO())
	return
}
