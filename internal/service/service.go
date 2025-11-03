package service

import (
	"context"
	"delayedNotifier/internal/entity"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/wb-go/wbf/rabbitmq"
)

type Storage interface {
	CreateNotify(ctx context.Context, notification entity.Notification) error
	GetNotify(ctx context.Context, id string) (*entity.Notification, error)
	DeleteNotify(ctx context.Context, id string) error
}

type Service struct {
	storage   Storage
	publisher *rabbitmq.Publisher
}

func NewService(storage Storage, publisher *rabbitmq.Publisher) *Service {
	return &Service{
		storage:   storage,
		publisher: publisher,
	}
}

func (s *Service) CreateNotify(ctx context.Context, req entity.Request) (string, error) {
	id := uuid.New().String()

	n := entity.Notification{
		ID:      id,
		Message: req.Message,
		UserID:  req.UserID,
		SendAt:  time.Now(),
		Status:  "pending",
	}

	err := s.storage.CreateNotify(ctx, n)
	if err != nil {
		return "", err
	}

	body, err := json.Marshal(n)
	if err != nil {
		return "", err
	}

	err = s.publisher.Publish(
		body,
		"notify",
		"application/json",
	)

	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *Service) GetNotify(ctx context.Context, id string) (*entity.Notification, error) {
	return s.storage.GetNotify(ctx, id)
}

func (s *Service) DeleteNotify(ctx context.Context, id string) error {
	return s.storage.DeleteNotify(ctx, id)
}
