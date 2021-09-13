package config

import (
	"log"
	"time"

	"github.com/duongnln96/building-microservices-golang/auth-service/tools/mongodb"
	"github.com/spf13/viper"
)

type AuthServiceConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"readtimeout"`
	WriteTimeout time.Duration `mapstructure:"writetimeout"`
	JWTSercet    string        `mapstructure:"jwt_sercet"`
	DBCollection string        `mapstructure:"db_collection"`
}

type AppConfig struct {
	MongoDBConfig mongodb.MongoDBConfig `mapstructure:"mongodb"`
	AuthConfig    AuthServiceConfig     `mapstructure:"authservice"`
}

var values AppConfig

func init() {
	config := viper.New()
	config.SetConfigName("config") // config file name
	config.AddConfigPath("./config/")
	// config.AutomaticEnv()

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
