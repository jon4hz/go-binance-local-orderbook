package handler

import (
	"context"
	"fmt"

	"github.com/adshao/go-binance/v2"
	"github.com/jon4hz/go-binance-local-orderbook/config"
)

var (
	lastUpdateID int64
	u            int64
	U            int64
	prev_u       int64
)

func HandleWebsocket(config config.Config) {
	wsDepthHandler := func(event *binance.WsDepthEvent) {
		U = event.FirstUpdateID
		u = event.UpdateID
		// first time
		if prev_u == 0 {
			prev_u = u
			snap, err := downloadSnapshot(config)
			if err != nil {
				panic("Error while downloading the snapshot")
			}
			lastUpdateID = snap.LastUpdateID
			fmt.Println(lastUpdateID)

		}
		fmt.Println(u, prev_u+1, U)
		prev_u = u

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
