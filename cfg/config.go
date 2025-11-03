package cfg

import (
	"github.com/wb-go/wbf/config"
	"github.com/wb-go/wbf/rabbitmq"
)

type AppConfig struct {
	Server serverConfig `mapstructure:"server"`
	Logger loggerConfig `mapstructure:"logger"`
	Rabbit RabbitConfig `mapstructure:"rabbit"`
}

type serverConfig struct {
	addr string
}

type loggerConfig struct {
	logLevel string
}

type RabbitConfig struct {
	Url      string                  `mapstructure:"url"`
	Consumer rabbitmq.ConsumerConfig `mapstructure:"consumer"`
}

func NewAppConfig() (*AppConfig, error) {
	cfg := config.New()

	cfg.EnableEnv("")
	config.ErrLoadConfigFile = cfg.LoadConfigFiles("./config.yaml")
	if config.ErrLoadConfigFile != nil {
		return nil, config.ErrLoadConfigFile
	}

	var appCfg AppConfig
	if err := cfg.Unmarshal(&appCfg); err != nil {
		return nil, err
	}

	return &appCfg, nil
}
