package app

import (
	"delayedNotifier/cfg"
	rabbit "delayedNotifier/internal/rabbitmq"
	"log"

	"github.com/wb-go/wbf/zlog"
)

func Run(path string) {
	// config
	config, err := cfg.NewAppConfig()
	if err != nil {
		log.Fatal(err)
	}
	// logger
	zlog.Init()

	// rabbit
	rmq, err := rabbit.New(config.Rabbit.Url, "notifications")
	if err != nil {
		zlog.Logger.Fatal()
	}

}
