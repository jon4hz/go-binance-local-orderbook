package binance

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/adshao/go-binance/v2"
	"github.com/jon4hz/go-binance-local-orderbook/alerting"
	"github.com/jon4hz/go-binance-local-orderbook/config"
	"github.com/jon4hz/go-binance-local-orderbook/database"
	"github.com/jon4hz/go-binance-local-orderbook/exchange"
)

func InitWebsocket(config *config.Config) {
	//var response database.DatabaseInsert
	wsDepthHandler := func(event *binance.WsDepthEvent) {
		exchange.Notified = false
		response := &database.BinanceDepthResponse{Response: event}
		exchange.BigU = event.FirstUpdateID
		exchange.SmallU = event.UpdateID
		// first time
		if exchange.Prev_u == 0 {
			// download snapshot
			if exchange.LastUpdateID == 0 {
				snap, err := downloadSnapshot(config)
				if err != nil {
					log.Println("Error while downloading the snapshot")
					return
				}
				response = &database.BinanceDepthResponse{Snapshot: snap}
				err = response.InsertIntoDatabase(config.Database.DBTableMarketName)
				if err != nil {
					log.Println(err)
					msg := alerting.AlertingMSG(fmt.Sprintf("🚨 Error coin: %s, couldn't insert data %s", config.Exchange.Market, err))
					go msg.TriggerAlert(config.Alerting)
					return
				}
				msg := alerting.AlertingMSG(fmt.Sprintf("💡 Info: Downloaded new snapshot for coin: %s", config.Exchange.Market))
				go msg.TriggerAlert(config.Alerting)
				log.Println("[binance][dbinsert] Inserted snapshot into db")
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
				log.Println("[binance][dbinsert] Inserted first event successfully")
			}
			return

		} else if exchange.BigU >= exchange.Prev_u+1 {
			if exchange.BigU > exchange.Prev_u+1 {
				log.Printf("Warning, U = %d and prev_u = %d", exchange.BigU, exchange.Prev_u)
				msg := alerting.AlertingMSG(fmt.Sprintf("⚠️ Warning: Orderbook could be out of sync for %s, U: %d, prev_u: %d", config.Exchange.Market, exchange.BigU, exchange.Prev_u))
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
			log.Printf("Websocket for %s crashed, spawning a new one.", sym)
			doneC, _, err := binance.WsDepthServe(sym, wsDepthHandler, errHandler)
			if err != nil {
				log.Printf("error registering symbol %s: %v", sym, err)
			}
			monitorWS(sym, doneC)
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

	}(config.Exchange.Market)
	wg.Wait()
}

func downloadSnapshot(config *config.Config) (res *binance.DepthResponse, err error) {
	client := binance.NewClient("", "")
	res, err = client.NewDepthService().Symbol(config.Exchange.Market).
		Limit(1000).
		Do(context.TODO())
	return
}
