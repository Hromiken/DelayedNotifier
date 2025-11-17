package cfg

import (
	"fmt"
	"os"
	"strings"

	"github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
)

type AppConfig struct {
	App      appInfo        `mapstructure:"app"`
	Server   serverConfig   `mapstructure:"server"`
	Logger   loggerConfig   `mapstructure:"logger"`
	Rabbit   RabbitConfig   `mapstructure:"rabbit"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Notifier NotifierConfig `mapstructure:"notifier"`
}

type appInfo struct {
	Name    string `yaml:"name" mapstructure:"name"`
	Version string `yaml:"version" mapstructure:"version"`
}
type serverConfig struct {
	Addr string `yaml:"addr" mapstructure:"addr"`
}

type loggerConfig struct {
	LogLevel string `yaml:"log_level" mapstructure:"log_level"`
}

type RabbitConfig struct {
	Url      string         `yaml:"url" mapstructure:"url"`
	Consumer ConsumerConfig `yaml:"consumer" mapstructure:"consumer"`
}

type ConsumerConfig struct {
	Queue     string        `yaml:"queue" mapstructure:"queue"`
	Consumer  string        `yaml:"consumer" mapstructure:"consumer"`
	AutoAck   bool          `yaml:"auto_ack" mapstructure:"auto_ack"`
	Exclusive bool          `yaml:"exclusive" mapstructure:"exclusive"`
	NoLocal   bool          `yaml:"no_local" mapstructure:"no_local"`
	NoWait    bool          `yaml:"no_wait" mapstructure:"no_wait"`
	Args      amqp091.Table `yaml:"args" mapstructure:"args"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr" mapstructure:"addr"`
	Password string `yaml:"password" mapstructure:"password"`
	DB       int    `yaml:"db" mapstructure:"db"`
}

type NotifierConfig struct {
	Type     string         `yaml:"type" mapstructure:"type"`
	Email    EmailConfig    `yaml:"email" mapstructure:"email"`
	Telegram TelegramConfig `yaml:"telegram" mapstructure:"telegram"`
}

type EmailConfig struct {
	Host     string `yaml:"host" mapstructure:"host"`
	Port     int    `yaml:"port" mapstructure:"port"`
	Username string `yaml:"username" mapstructure:"username"`
	Password string `yaml:"password" mapstructure:"password"`
	From     string `yaml:"from" mapstructure:"from"`
}

type TelegramConfig struct {
	Token         string `yaml:"token" mapstructure:"token"`
	DefaultChatID string `yaml:"default_chat_id" mapstructure:"default_chat_id"`
}

func NewAppConfig() (*AppConfig, error) {
	viper.SetConfigFile("./cfg/config.yaml")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// связываем env-переменные с полями конфига
	_ = viper.BindEnv("rabbit.url", "RABBITMQ_URL")
	_ = viper.BindEnv("redis.addr", "REDIS_ADDR")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("CONFIG ERROR:", err)
		return nil, err
	}

	fmt.Println("Used config:", viper.ConfigFileUsed())

	var appCfg AppConfig
	err = viper.Unmarshal(&appCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if appCfg.Rabbit.Url == "" {
		envURL := os.Getenv("RABBITMQ_URL")
		if envURL != "" {
			appCfg.Rabbit.Url = envURL
		}
	}

	return &appCfg, nil
}
