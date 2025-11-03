package rabbit

import (
	"log"
	"time"

	"github.com/wb-go/wbf/rabbitmq"
)

type RMQ struct {
	Conn      *rabbitmq.Connection
	Channel   *rabbitmq.Channel
	Publisher *rabbitmq.Publisher
}

func New(url string, exchangeName string) (*RMQ, error) {
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

	pub := rabbitmq.NewPublisher(ch, exchangeName)

	log.Println("RabbitMQ connected")

	return &RMQ{
		Conn:      conn,
		Channel:   ch,
		Publisher: pub,
	}, nil
}
