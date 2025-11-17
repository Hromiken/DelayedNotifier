package sender

import (
	"context"
	"delayedNotifier/internal/entity"
	"log"
)

type MockSender struct{}

func NewMockSender() *MockSender {
	return &MockSender{}
}

func (s *MockSender) Send(_ context.Context, n entity.Notification) error {
	log.Printf("MockSender: sending notification id=%s to user=%s message=%s",
		n.ID, n.UserID, n.Message)
	return nil
}
