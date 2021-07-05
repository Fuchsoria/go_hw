package main

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Logger LoggerConf `json:"logger"`
	AMPQ   AMPQConf   `json:"ampq"`
}

type LoggerConf struct {
	Level string `json:"level"`
	File  string `json:"file"`
}

type AMPQConf struct {
	URI  string `json:"uri"`
	Name string `json:"name"`
}

func NewConfig() (Config, error) {
	viper.SetConfigFile(configFile)

	if err := viper.ReadInConfig(); err != nil { // Handle errors reading the config file
		return Config{}, fmt.Errorf("fatal error config file: %w", err)
	}

	return Config{
			LoggerConf{Level: viper.GetString("logger.level"), File: viper.GetString("logger.file")},
			AMPQConf{URI: viper.GetString("ampq.uri"), Name: viper.GetString("ampq.name")},
		},
		nil
}
