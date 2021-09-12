package mongodb

import (
	"fmt"
	"time"
)

type MongoDBConfig struct {
	Host                string        `mapstructure:"host"`
	Port                int           `mapstructure:"port"`
	Username            string        `mapstructure:"username"`
	Password            string        `mapstructure:"password"`
	DBName              string        `mapstructure:"dbname"`
	Timeout             time.Duration `mapstructure:"querytimeout"`
	HealthCheckInterval time.Duration `mapstructure:"healthcheck_interval"`
}

func (mc *MongoDBConfig) GetInfo() string {
	return fmt.Sprintf(
		"mongodb://%s:%s@%s:%d/",
		mc.Username, mc.Password, mc.Host, mc.Port,
	)
}
