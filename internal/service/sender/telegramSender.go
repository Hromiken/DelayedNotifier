package sender

import (
	"bytes"
	"context"
	"delayedNotifier/cfg"
	"delayedNotifier/internal/entity"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type TelegramSender struct {
	cfg cfg.TelegramConfig
}

func NewTelegramSender(c cfg.TelegramConfig) *TelegramSender {
	return &TelegramSender{cfg: c}
}

func (s *TelegramSender) Send(ctx context.Context, n entity.Notification) error {
	chatID := n.UserID
	if chatID == "" {
		chatID = s.cfg.DefaultChatID
	}

	req := map[string]interface{}{
		"chat_id": chatID,
		"text":    n.Message,
	}

	body, _ := json.Marshal(req)

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", s.cfg.Token)

	httpClient := &http.Client{Timeout: 10 * time.Second}

	r, err := httpClient.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		return fmt.Errorf("telegram returned status %d", r.StatusCode)
	}

	return nil
}
