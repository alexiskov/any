package confreader

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type (
	Configs struct {
		DMS  *DataBase
		Tbot *TbotData
	}
	TbotData struct {
		API string `env:"TGBOT_APIKEY"`
	}

	DataBase struct {
		Host     string `env:"DB_HOST"`
		Port     int    `env:"DB_PORT"`
		DBname   string `env:"DB_NAME"`
		User     string `env:"DB_USER"`
		Password string `env:"DB_PASSWORD"`
		SSLmode  string `env:"DB_SSLMODE"`
	}
)

func LoadConfig() (c Configs, err error) {
	if err = godotenv.Load(); err != nil {
		err = fmt.Errorf("config loading -> env-file loading error: %w", err)
		return
	}
	c = Configs{&DataBase{Host: os.Getenv("DB_HOST"), DBname: os.Getenv("DB_NAME"), User: os.Getenv("DB_USER"), Password: os.Getenv("DB_PASSWORD"), SSLmode: os.Getenv("DB_SSLMODE")}, &TbotData{API: os.Getenv("TGBOT_APIKEY")}}

	if dbport, err := strconv.Atoi(os.Getenv("DB_PORT")); err != nil {
		err = fmt.Errorf("config field DB_PORT parse error: %w", err)
		return Configs{}, err
	} else {
		c.DMS.Port = dbport
	}

	return
}
