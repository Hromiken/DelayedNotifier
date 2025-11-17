package service

import (
	"context"
	"delayedNotifier/internal/entity"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/wb-go/wbf/rabbitmq"
)

type Sender interface {
	Send(ctx context.Context, notification entity.Notification) error
}

type Worker struct {
	storage        Storage
	defaultSender  Sender
	emailSender    Sender
	telegramSender Sender
	consumer       *rabbitmq.Consumer
	publisher      *rabbitmq.Publisher
	msgChan        chan []byte
}

func NewWorker(storage Storage, defaultSender, emailSender, telegramSender Sender,
	consumer *rabbitmq.Consumer,
	publisher *rabbitmq.Publisher,
) *Worker {
	return &Worker{
		storage:        storage,
		defaultSender:  defaultSender,
		emailSender:    emailSender,
		telegramSender: telegramSender,
		consumer:       consumer,
		publisher:      publisher,
		msgChan:        make(chan []byte, 100),
	}
}

func (w *Worker) chooseSender(n *entity.Notification) Sender {
	if n.Sender == "" {
		log.Printf("worker: empty sender for id=%s, using default(mock)", n.ID)
		return w.defaultSender
	}

	switch n.Sender {
	case "email":
		return w.emailSender
	case "telegram":
		return w.telegramSender
	case "mock":
		return w.defaultSender // default уже MultiSender(maybe mock)
	default:
		log.Printf("worker: unknown sender '%s' id=%s -> mock", n.Sender, n.ID)
		return w.defaultSender
	}
}

func (w *Worker) Start(ctx context.Context) error {
	go func() {
		if err := w.consumer.Consume(w.msgChan); err != nil {
			log.Printf("worker: consumer stopped with error: %v", err)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			log.Println("worker: context cancelled, stopping worker loop")
			return ctx.Err()

		case body, ok := <-w.msgChan:
			if !ok {
				log.Println("worker: messages channel closed")
				return nil
			}
			go w.processMessage(ctx, body)
		}
	}
}

func (w *Worker) processMessage(parentCtx context.Context, body []byte) {
	var payload map[string]string
	if err := json.Unmarshal(body, &payload); err != nil {
		log.Printf("worker: bad message body: %v", err)
		return
	}

	id := payload["id"]
	if id == "" {
		log.Println("worker: message without id")
		return
	}

	if err := w.handleNotification(parentCtx, id); err != nil {
		log.Printf("worker: handling id=%s finished with err: %v", id, err)
	}
}

func (w *Worker) handleNotification(ctx context.Context, id string) error {
	n, err := w.storage.GetNotify(ctx, id)
	if err != nil {
		log.Printf("worker: get notify failed id=%s err=%v", id, err)
		go w.scheduleRepublish(ctx, id, 1)
		return err
	}

	if n.Status == "cancelled" {
		return nil
	}

	now := time.Now()
	if now.Before(n.SendAt) {
		n.Status = StatusScheduled
		_ = w.storage.CreateNotify(ctx, *n)

		wait := time.Until(n.SendAt)
		timer := time.NewTimer(wait)
		select {
		case <-timer.C:
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		}
	}

	n, err = w.storage.GetNotify(ctx, id)
	if err != nil {
		go w.scheduleRepublish(ctx, id, 1)
		return err
	}

	if n.Status == "cancelled" {
		return nil
	}

	sender := w.chooseSender(n)
	if sender == nil {
		return errors.New("no sender found")
	}

	err = sender.Send(ctx, *n)
	if err != nil {
		log.Printf("worker: send failed id=%s err=%v", id, err)
		n.Retry++
		n.Status = StatusFailed
		_ = w.storage.CreateNotify(ctx, *n)
		go w.scheduleRepublish(ctx, id, n.Retry)
		return err
	}

	n.Status = StatusSent
	n.Retry = 0
	_ = w.storage.CreateNotify(ctx, *n)
	log.Printf("worker: sent id=%s", id)
	return nil
}

func (w *Worker) scheduleRepublish(parentCtx context.Context, id string, retryCount int) {
	if retryCount > 6 {
		ctx, cancel := context.WithTimeout(parentCtx, 5*time.Second)
		defer cancel()
		if n, err := w.storage.GetNotify(ctx, id); err == nil {
			n.Status = StatusFailed
			n.Retry = retryCount
			_ = w.storage.CreateNotify(ctx, *n)
		}
		return
	}

	delay := time.Duration(1<<retryCount) * time.Second
	timer := time.NewTimer(delay)

	select {
	case <-timer.C:
		body, _ := json.Marshal(map[string]string{"id": id})
		if err := w.publisher.Publish(body, "notify", "application/json"); err != nil {
			go w.scheduleRepublish(parentCtx, id, retryCount+1)
			return
		}
	case <-parentCtx.Done():
		timer.Stop()
		return
	}
}
