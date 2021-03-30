package handler

import (
	"context"
	"fmt"

	"github.com/adshao/go-binance/v2"
	"github.com/jon4hz/go-binance-local-orderbook/config"
)

var (
	LastUpdateID int64
	SmallU       int64
	BigU         int64
	Prev_u       int64
)

func HandleWebsocket(config config.Config) {
	wsDepthHandler := func(event *binance.WsDepthEvent) {
		BigU = event.FirstUpdateID
		SmallU = event.UpdateID
		// first time
		if Prev_u == 0 {
			Prev_u = SmallU
			snap, err := downloadSnapshot(config)
			if err != nil {
				panic("Error while downloading the snapshot")
			}
			LastUpdateID = snap.LastUpdateID
			fmt.Println(LastUpdateID)

		}
		fmt.Println(SmallU, Prev_u+1, BigU)
		Prev_u = SmallU

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
