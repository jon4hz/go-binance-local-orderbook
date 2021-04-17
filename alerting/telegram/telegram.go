package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jon4hz/go-binance-local-orderbook/config"
)

func TriggerTelegramAlert(cfg *config.Config, msg interface{}) error {
	client := &http.Client{}
	data := map[string]interface{}{"chat_id": cfg.Alerting.Telegram.Chat, "text": msg}
	jsonStr, err := json.Marshal(data)
	if err != nil {
		return err
	}
	request_url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", cfg.Alerting.Telegram.Token)
	req, err := http.NewRequest(http.MethodPost, request_url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	switch res.StatusCode {
	case 200:
		return nil
	default:
		return fmt.Errorf("%s", res.Status)
	}
}
