package main

//nolint:depguard
import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Logger     LoggerConf
	SQL        SQLConf
	HTTPServer HTTPServerConf
	GRPCServer GRPCServerConf
}

type LoggerConf struct {
	Level string `mapstructure:"level" default:"INFO"`
}

type SQLConf struct {
	Username       string `mapstructure:"userName"`
	Password       string `mapstructure:"password"`
	Host           string `mapstructure:"Host"`
	Port           string `mapstructure:"Port"`
	Database       string `mapstructure:"database"`
	MigrationsPath string `mapstructure:"migrationsPath"`
}

type HTTPServerConf struct {
	Host string `mapstructure:"host" default:"0.0.0.0"`
	Port string `mapstructure:"port" default:"8080"`
}

type GRPCServerConf struct {
	Port string `mapstructure:"port" default:"50051"`
}

func NewConfig(path string) (Config, error) {
	var conf Config
	err := viper.BindEnv("SQL.Host", "sqlHost")
	if err != nil {
		return conf, err
	}

	err = viper.BindEnv("SQL.Port", "sqlPort")
	if err != nil {
		return conf, err
	}

	err = viper.BindEnv("SQL.Database", "sqlDatabase")
	if err != nil {
		return conf, err
	}

	viper.SetDefault("SQL.Username", "postgres")
	viper.SetDefault("SQL.Password", "password")
	viper.SetDefault("SQL.Host", "0.0.0.0")
	viper.SetDefault("SQL.Port", "5435")
	viper.SetDefault("SQL.Database", "backend")
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return conf, fmt.Errorf("error while reading config file: %w", err)
	}

	if err := viper.Unmarshal(&conf); err != nil {
		return conf, fmt.Errorf("error while unmarshaling config: %w", err)
	}

	return conf, nil
}
