package rabbit

import (
	"fmt"
	"log"
	"time"

	"github.com/wb-go/wbf/rabbitmq"
)

type RMQ struct {
	Conn      *rabbitmq.Connection
	Channel   *rabbitmq.Channel
	Publisher *rabbitmq.Publisher
}

func NewPublisher(url string, exchangeName string) (*RMQ, error) {
	conn, err := rabbitmq.Connect(url, 5, time.Second)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	exchange := rabbitmq.NewExchange(exchangeName, "direct")
	exchange.Durable = true

	err = exchange.BindToChannel(ch)
	if err != nil {
		return nil, err
	}

	log.Println("RabbitMQ connected")

	queue, err := ch.QueueDeclare(
		"notifications",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("queue declare failed: %w", err)
	}

	err = ch.QueueBind(
		queue.Name,
		"notify",     // routing key
		exchangeName, // exchange
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("queue bind failed: %w", err)
	}

	pub := rabbitmq.NewPublisher(ch, exchangeName)

	return &RMQ{
		Conn:      conn,
		Channel:   ch,
		Publisher: pub,
	}, nil
}
