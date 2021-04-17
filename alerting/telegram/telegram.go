package telegram

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/jon4hz/go-binance-local-orderbook/config"
)

func TriggerTelegramAlert(cfg *config.Config, msg interface{}) error {
	client := &http.Client{}
	jsonStr := []byte(fmt.Sprintf(`{"chat_id": %d, "text": "%s"}`, cfg.Alerting.Telegram.Chat, msg))
	request_url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", cfg.Alerting.Telegram.Token)
	req, err := http.NewRequest(http.MethodPost, request_url, bytes.NewBuffer(jsonStr))
	fmt.Println(request_url)
	fmt.Println(string(jsonStr))
	if err != nil {
		return err
	}
	req.Header.Set("content-type", "application-json")
	res, err := client.Do(req)
	if err != nil {
		return err
		//return fmt.Errorf("%s %s", res.Status, err)
	}
	switch res.Status {
	case "200":
		return nil
	default:
		return fmt.Errorf("%s", res.Status)
	}
}
