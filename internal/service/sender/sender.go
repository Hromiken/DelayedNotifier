package sender

import (
	"context"
	"delayedNotifier/cfg"
	"delayedNotifier/internal/entity"
	"fmt"
)

type Sender interface {
	Send(ctx context.Context, n entity.Notification) error
}

type MultiSender struct {
	Email    Sender
	Telegram Sender
	Mock     Sender
}

func NewMultiSender(email Sender, telegram Sender, mock Sender) *MultiSender {
	return &MultiSender{
		Email:    email,
		Telegram: telegram,
		Mock:     mock,
	}
}

func NewSender(cfg cfg.NotifierConfig) (Sender, error) {
	switch cfg.Type {
	case "multi":
		email := NewEmailSender(cfg.Email)
		telegram := NewTelegramSender(cfg.Telegram)
		mock := NewMockSender()
		return NewMultiSender(email, telegram, mock), nil

	case "email":
		return NewEmailSender(cfg.Email), nil

	case "telegram":
		return NewTelegramSender(cfg.Telegram), nil

	case "mock":
		return NewMockSender(), nil

	default:
		return nil, fmt.Errorf("unknown sender type: %s", cfg.Type)
	}
}

func (m *MultiSender) Send(ctx context.Context, n entity.Notification) error {
	switch n.Channel {
	case "email":
		return m.Email.Send(ctx, n)
	case "telegram":
		return m.Telegram.Send(ctx, n)
	case "mock":
		return m.Mock.Send(ctx, n)
	case "all":
		if err := m.Email.Send(ctx, n); err != nil {
			return err
		}
		if err := m.Telegram.Send(ctx, n); err != nil {
			return err
		}
		return m.Mock.Send(ctx, n)
	default:
		return fmt.Errorf("unknown channel: %s", n.Channel)
	}
}
