package entity

import "time"

// Notification - сущность уведомления
type Notification struct {
	ID      string    `json:"id"`
	Message string    `json:"message"`
	UserID  string    `json:"user_id"`
	SendAt  time.Time `json:"send_at"`
	Status  string    `json:"status"`
	Retry   int       `json:"retry"`
	Sender  string    `json:"sender"`
	Channel string
}

type Request struct {
	Message string    `json:"message"`
	UserID  string    `json:"user_id"`
	SendAt  time.Time `json:"send_at"`
	Sender  string    `json:"sender"`
}
