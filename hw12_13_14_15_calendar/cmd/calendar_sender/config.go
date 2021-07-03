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
	uri  string
	name string
}

func NewConfig() (Config, error) {
	viper.SetConfigFile(configFile)

	if err := viper.ReadInConfig(); err != nil { // Handle errors reading the config file
		return Config{}, fmt.Errorf("fatal error config file: %w", err)
	}

	return Config{
			LoggerConf{viper.GetString("logger.level"), viper.GetString("logger.file")},
			AMPQConf{viper.GetString("ampq.uri"), viper.GetString("ampq.name")},
		},
		nil
}
