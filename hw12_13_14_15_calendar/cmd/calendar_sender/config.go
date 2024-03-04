package main

//nolint:depguard
import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Logger LoggerConf
	MB     MBConf
}

type LoggerConf struct {
	Level string `mapstructure:"level" default:"INFO"`
}

type MBConf struct {
	Username         string `mapstructure:"username"`
	Password         string `mapstructure:"password"`
	Host             string `mapstructure:"host"`
	Port             string `mapstructure:"port"`
	Protocol         string `mapstructure:"protocol"`
	ExchangeName     string `mapstructure:"exchangeName"`
	ExchangeType     string `mapstructure:"exchangeType"`
	QueueName        string `mapstructure:"queueName"`
	ConfirmQueueName string `mapstructure:"confirmQueueName"`
	RouteKey         string `mapstructure:"routeKey"`
	ClientTag        string `mapstructure:"clientTag"`
}

func NewConfig(path string) (Config, error) {
	var conf Config
	err := viper.BindEnv("MB.Host", "mbHost")
	if err != nil {
		return conf, err
	}

	err = viper.BindEnv("MB.Port", "mbPort")
	if err != nil {
		return conf, err
	}

	err = viper.BindEnv("MB.Username", "mbUsername")
	if err != nil {
		return conf, err
	}

	err = viper.BindEnv("MB.Password", "mbPassword")
	if err != nil {
		return conf, err
	}

	viper.SetDefault("MB.Username", "rabbit")
	viper.SetDefault("MB.Password", "password")
	viper.SetDefault("MB.Host", "0.0.0.0")
	viper.SetDefault("MB.Port", "5672")
	viper.SetDefault("MB.Protocol", "amqp")
	viper.SetDefault("MB.ExchangeName", "test-exchange")
	viper.SetDefault("MB.ExchangeType", "topic")
	viper.SetDefault("MB.QueueName", "test-queue")
	viper.SetDefault("MB.ConfirmQueueName", "ToNotification")
	viper.SetDefault("MB.RouteKey", "test-route")
	viper.SetDefault("MB.ClientTag", "test-client")

	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return conf, fmt.Errorf("error while reading config file: %w", err)
	}

	if err := viper.Unmarshal(&conf); err != nil {
		return conf, fmt.Errorf("error while unmarshaling config: %w", err)
	}

	return conf, nil
}
