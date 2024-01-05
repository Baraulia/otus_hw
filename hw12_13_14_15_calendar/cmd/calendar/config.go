package main

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Logger LoggerConf
	// TODO
}

type LoggerConf struct {
	Level string
	// TODO
}

func NewConfig(path string) (Config, error) {
	var conf Config
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return conf, fmt.Errorf("error while reading cinfig file: %s", err)
	}

	if err := viper.Unmarshal(&conf); err != nil {
		return conf, fmt.Errorf("error while unmarshaling config: %s", err)
	}

	return conf, nil
}
