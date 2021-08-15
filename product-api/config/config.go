package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"readtimeout"`
	WriteTimeout time.Duration `mapstructure:"writetimeout"`
}

type CurrencyService struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type AppConfig struct {
	Server      ServerConfig    `mapstructure:"productapi"`
	CurrService CurrencyService `mapstructure:"currency"`
}

var values AppConfig

func init() {
	config := viper.New()
	config.SetConfigName("config") // config file name
	config.AddConfigPath("./config/")
	// config.AutomaticEnv()

	err := config.ReadInConfig()
	if err != nil {
		log.Fatalf("error read config / %s", err)
	}

	if err := config.ReadInConfig(); err != nil {
		log.Fatalf("Error while reading config: %s", err)
	}

	if err := config.Unmarshal(&values); err != nil {
		log.Fatalf("Error while parsing config: %s", err)
	}
}

func GetConfig() *AppConfig {
	return &values
}
