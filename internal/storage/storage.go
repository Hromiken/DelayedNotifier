package storage

import (
	"context"
	"delayedNotifier/internal/entity"
	"encoding/json"
	"fmt"
	"log"

	"github.com/wb-go/wbf/redis"
)

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(client *redis.Client) *RedisStorage {
	return &RedisStorage{client: client}
}

// CreateNotify - создает уведомление
func (r *RedisStorage) CreateNotify(ctx context.Context, notification entity.Notification) error {
	marshal, err := json.Marshal(notification)
	if err != nil {
		log.Println(err)
		return err
	}
	key := fmt.Sprintf("notify:%s:payload", notification.ID)

	err = r.client.Set(ctx, key, string(marshal))
	if err != nil {
		return err
	}
	return nil
}

// GetNotify - получить уведомление по id
func (r *RedisStorage) GetNotify(ctx context.Context, id string) (*entity.Notification, error) {
	key := fmt.Sprintf("notify:%s:payload", id)

	val, err := r.client.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	var n entity.Notification
	err = json.Unmarshal([]byte(val), &n)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &n, nil
}

// DeleteNotify - удалить уведомление по id
func (r *RedisStorage) DeleteNotify(ctx context.Context, id string) error {
	key := fmt.Sprintf("notify:%s:payload", id)

	err := r.client.Del(ctx, key)
	if err != nil {
		return err
	}
	return nil
}
