package confreader

import (
	"fmt"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type (
	Configs struct {
		HTTPserver struct {
			Port     int    `env:"HTTP_SERVER_PORT"`
			TimeZone string `env:"HTTP_SERVER_TIMEZONE"`
		}
		DataBase struct {
			Host     string `env:"DB_HOST"`
			Port     int    `env:"DB_PORT"`
			DBname   string `env:"DB_NAME"`
			User     string `env:"DB_USER"`
			Password string `env:"DB_PASSWORD"`
			SSL      bool   `env:"DB_SSLMODE"`
		}
	}
)

func LoadConfig() (c Configs, err error) {
	if err = godotenv.Load(); err != nil {
		err = fmt.Errorf("config loading -> env-file loading error: %w", err)
		return
	}
	if err = env.Parse(&c); err != nil {
		err = fmt.Errorf("config loading -> env-filedata parsing error: %w", err)
		return
	}
	return
}
