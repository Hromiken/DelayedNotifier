package service

import (
	"context"
	"delayedNotifier/internal/entity"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/wb-go/wbf/rabbitmq"
)

const (
	StatusPending   = "pending"
	StatusScheduled = "scheduled"
	StatusSent      = "sent"
	StatusFailed    = "failed"
)

type Storage interface {
	CreateNotify(ctx context.Context, notification entity.Notification) error
	GetNotify(ctx context.Context, id string) (*entity.Notification, error)
	DeleteNotify(ctx context.Context, id string) error
}

// NotifierService - delayed notifier service.
type NotifierService struct {
	storage   Storage
	publisher *rabbitmq.Publisher
}

// NewNotifierService - конструктор NotifierService.
func NewNotifierService(storage Storage, publisher *rabbitmq.Publisher) *NotifierService {
	return &NotifierService{
		storage:   storage,
		publisher: publisher,
	}
}

// CreateNotify - метод NotifierService для создания уведомления и отправления его ID в очередь RabbitMQ.
func (s *NotifierService) CreateNotify(ctx context.Context, req entity.Request) (string, error) {
	id := uuid.New().String()

	n := entity.Notification{
		ID:      id,
		Message: req.Message,
		UserID:  req.UserID,
		SendAt:  req.SendAt,
		Status:  StatusPending,
		Retry:   0,
		Sender:  req.Sender,
	}

	err := s.storage.CreateNotify(ctx, n)
	if err != nil {
		return "", err
	}

	body, err := json.Marshal(map[string]string{"id": id})
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

// GetNotify - метод NotifierService для получения уведомления по id.
func (s *NotifierService) GetNotify(ctx context.Context, id string) (*entity.Notification, error) {
	return s.storage.GetNotify(ctx, id)
}

// DeleteNotify - метод NotifierService для удаления уведомления по id.
func (s *NotifierService) DeleteNotify(ctx context.Context, id string) error {
	return s.storage.DeleteNotify(ctx, id)
}
