package app

import (
	"context"
	config2 "delayedNotifier/cfg"
	"delayedNotifier/internal/handler"
	rabbit "delayedNotifier/internal/rabbitmq"
	"delayedNotifier/internal/service"
	"delayedNotifier/internal/service/sender"
	"delayedNotifier/internal/storage"
	"log"
	"net/http"

	"github.com/wb-go/wbf/rabbitmq"
	"github.com/wb-go/wbf/zlog"
)

func Run() {
	// config
	config, err := config2.NewAppConfig()
	if err != nil {
		log.Fatal(err)
	}

	// logger
	zlog.Init()

	// rabbit
	rmq, err := rabbit.NewPublisher(config.Rabbit.Url, "notifications")
	if err != nil {
		zlog.Logger.Fatal().Msg(err.Error())
	}

	// storage
	st := storage.NewRedisStorage(config)

	// service
	svc := service.NewNotifierService(st, rmq.Publisher)

	// consumer
	cfgConsumer := &rabbitmq.ConsumerConfig{
		Queue:     config.Rabbit.Consumer.Queue,
		Consumer:  config.Rabbit.Consumer.Consumer,
		AutoAck:   config.Rabbit.Consumer.AutoAck,
		Exclusive: config.Rabbit.Consumer.Exclusive,
		NoLocal:   config.Rabbit.Consumer.NoLocal,
		NoWait:    config.Rabbit.Consumer.NoWait,
		Args:      config.Rabbit.Consumer.Args,
	}
	consumerRab := rabbitmq.NewConsumer(rmq.Channel, cfgConsumer)

	// sender
	emailSender := sender.NewEmailSender(config.Notifier.Email)
	telegramSender := sender.NewTelegramSender(config.Notifier.Telegram)

	// default — мульт-sender
	defaultSender := sender.NewMultiSender(emailSender, telegramSender, sender.NewMockSender())

	// worker
	worker := service.NewWorker(
		st,
		defaultSender,
		emailSender,
		telegramSender,
		consumerRab,
		rmq.Publisher,
	)

	// ctx
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go worker.Start(ctx)

	hand := handler.NewNotifyHandler(svc)
	router := handler.NewRouter(hand)

	http.ListenAndServe(":8080", router)
}
