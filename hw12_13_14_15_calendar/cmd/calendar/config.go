package main

import (
	"fmt"

	"github.com/spf13/viper"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger LoggerConf `json:"logger"`
	DB     DBConf     `json:"db"`
	HTTP   HTTPConf   `json:"http"`
}

type LoggerConf struct {
	Level string `json:"level"`
	File  string `json:"file"`
}

type DBConf struct {
	Method           string `json:"method"`
	ConnectionString string `json:"connection_string"`
}

type HTTPConf struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

func NewConfig() (Config, error) {
	viper.SetConfigFile(configFile)

	if err := viper.ReadInConfig(); err != nil { // Handle errors reading the config file
		return Config{}, fmt.Errorf("fatal error config file: %w", err)
	}

	return Config{
		LoggerConf{viper.GetString("logger.level"), viper.GetString("logger.file")},
		DBConf{viper.GetString("db.method"), viper.GetString("db.connection_string")},
		HTTPConf{viper.GetString("http.host"), viper.GetString("http.port")},
	}, nil
}
