package config

import (
	"log"

	"github.com/spf13/viper"
)

type CurrencyConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

var values CurrencyConfig

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

func GetConfig() *CurrencyConfig {
	return &values
}
