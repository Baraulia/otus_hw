package main

//nolint:depguard
import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Logger LoggerConf
	SQL    SQLConf
	MB     MBConf
}

type LoggerConf struct {
	Level string `mapstructure:"level" default:"INFO"`
}

type SQLConf struct {
	Username       string `mapstructure:"userName"`
	Password       string `mapstructure:"password"`
	Host           string `mapstructure:"host"`
	Port           string `mapstructure:"port"`
	Database       string `mapstructure:"database"`
	MigrationsPath string `mapstructure:"migrationsPath"`
}

type MBConf struct {
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	Host         string `mapstructure:"host"`
	Port         string `mapstructure:"port"`
	Protocol     string `mapstructure:"protocol"`
	ExchangeName string `mapstructure:"exchangeName"`
	ExchangeType string `mapstructure:"exchangeType"`
	QueueName    string `mapstructure:"queueName"`
	RouteKey     string `mapstructure:"routeKey"`
}

func NewConfig(path string) (Config, error) {
	var conf Config
	viper.SetDefault("SQL.Username", "postgres")
	viper.SetDefault("SQL.Password", "password")
	viper.SetDefault("SQL.Host", "0.0.0.0")
	viper.SetDefault("SQL.Port", "5435")
	viper.SetDefault("SQL.Database", "backend")

	viper.SetDefault("MB.Username", "rabbit")
	viper.SetDefault("MB.Password", "password")
	viper.SetDefault("MB.Host", "0.0.0.0")
	viper.SetDefault("MB.Port", "5672")
	viper.SetDefault("MB.Protocol", "amqp")
	viper.SetDefault("MB.ExchangeName", "test-exchange")
	viper.SetDefault("MB.ExchangeType", "topic")
	viper.SetDefault("MB.QueueName", "test-queue")
	viper.SetDefault("MB.RouteKey", "test-route")

	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return conf, fmt.Errorf("error while reading config file: %w", err)
	}

	if err := viper.Unmarshal(&conf); err != nil {
		return conf, fmt.Errorf("error while unmarshaling config: %w", err)
	}

	return conf, nil
}
