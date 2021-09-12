package postgresql

import (
	"fmt"
	"time"
)

type PsqlConfig struct {
	Host                string        `mapstructure:"host"`
	Port                int           `mapstructure:"port"`
	UserName            string        `mapstructure:"username"`
	Password            string        `mapstructure:"password"`
	DBName              string        `mapstructure:"dbname"`
	TimeOut             time.Duration `mapstructure:"querytimeout"`
	HealthCheckInterval time.Duration `mapstructure:"healthcheck_interval"`
}

func (c *PsqlConfig) GetPsqlInfo() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.UserName, c.Password, c.DBName,
	)
}
