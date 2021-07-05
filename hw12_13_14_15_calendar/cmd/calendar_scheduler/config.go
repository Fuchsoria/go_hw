package main

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Scheduler SchedulerConfig `json:"scheduler"`
	DB        DBConf          `json:"db"`
	Logger    LoggerConf      `json:"logger"`
	AMPQ      AMPQConf        `json:"ampq"`
}

type SchedulerConfig struct {
	RecheckDelaySeconds int64 `json:"recheck_delay_seconds"`
}

type DBConf struct {
	ConnectionString string `json:"connection_string"`
}

type AMPQConf struct {
	URI  string `json:"uri"`
	Name string `json:"name"`
}

type LoggerConf struct {
	Level string `json:"level"`
	File  string `json:"file"`
}

func NewConfig() (Config, error) {
	viper.SetConfigFile(configFile)

	if err := viper.ReadInConfig(); err != nil { // Handle errors reading the config file
		return Config{}, fmt.Errorf("fatal error config file: %w", err)
	}

	return Config{
		SchedulerConfig{RecheckDelaySeconds: viper.GetInt64("scheduler.recheck_delay_seconds")},
		DBConf{ConnectionString: viper.GetString("db.connection_string")},
		LoggerConf{Level: viper.GetString("logger.level"), File: viper.GetString("logger.file")},
		AMPQConf{URI: viper.GetString("ampq.uri"), Name: viper.GetString("ampq.name")},
	}, nil
}
