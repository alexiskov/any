package confreader

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type (
	Configs struct {
		WebServer *HTTPserver
		DMS       *DataBase
		Tbot      *TbotData
	}
	TbotData struct {
		API string `env:"TGBOT_APIKEY"`
	}
	HTTPserver struct {
		Port     int    `env:"HTTPSERVER_PORT"`
		TimeZone string `env:"HTTPSERVER_TIMEZONE"`
	}
	DataBase struct {
		Host     string `env:"DB_HOST"`
		Port     int    `env:"DB_PORT"`
		DBname   string `env:"DB_NAME"`
		User     string `env:"DB_USER"`
		Password string `env:"DB_PASSWORD"`
		SSLmode  bool   `env:"DB_SSLMODE"`
	}
)

func LoadConfig() (c Configs, err error) {
	if err = godotenv.Load(); err != nil {
		err = fmt.Errorf("config loading -> env-file loading error: %w", err)
		return
	}
	c = Configs{&HTTPserver{TimeZone: os.Getenv("HTTPSERVER_TIMEZONE")}, &DataBase{Host: os.Getenv("DB_HOST"), DBname: os.Getenv("DB_NAME"), User: os.Getenv("DB_USER"), Password: os.Getenv("DB_PASSWORD")}, &TbotData{API: os.Getenv("TGBOT_APIKEY")}}

	if htpsport, err := strconv.Atoi(os.Getenv("HTTPSERVER_PORT")); err != nil {
		err = fmt.Errorf("config field HTTPSERVER_PORT parse error: %w", err)
		return Configs{}, err
	} else {
		c.WebServer.Port = htpsport
	}

	if dbport, err := strconv.Atoi(os.Getenv("DB_PORT")); err != nil {
		err = fmt.Errorf("config field DB_PORT parse error: %w", err)
		return Configs{}, err
	} else {
		c.DMS.Port = dbport
	}

	if sslmode, err := strconv.ParseBool(os.Getenv("DB_SSLMODE")); err != nil {
		err = fmt.Errorf("config field DB_SSLMODE parse error: %w", err)
		return Configs{}, err
	} else {
		c.DMS.SSLmode = sslmode
	}

	return
}
