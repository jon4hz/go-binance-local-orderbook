package binance

import (
	"context"
	"fmt"

	"github.com/adshao/go-binance/v2"
	"github.com/jon4hz/go-binance-local-orderbook/config"
	"github.com/jon4hz/go-binance-local-orderbook/handler"
)

func HandleWebsocket(config *config.Config) {
	wsDepthHandler := func(event *binance.WsDepthEvent) {
		handler.BigU = event.FirstUpdateID
		handler.SmallU = event.UpdateID
		// first time
		if handler.Prev_u == 0 {
			handler.Prev_u = handler.SmallU
			snap, err := downloadSnapshot(*config)
			if err != nil {
				panic("Error while downloading the snapshot")
			}
			handler.LastUpdateID = snap.LastUpdateID
			fmt.Println(handler.LastUpdateID)

		}
		fmt.Println(handler.SmallU, handler.Prev_u+1, handler.BigU)
		handler.Prev_u = handler.SmallU

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
