package main

import (
	"os"

	"vacancydealer/bd"
	"vacancydealer/confreader"
	"vacancydealer/hh"
	"vacancydealer/logger"
	"vacancydealer/telebot"
)

func main() {
	logger.InitInfoTextlog(os.Stdout)
	logger.Info("logger status is Run...")

	logger.InitErrorJSONlog(os.Stdout)
	logger.Info("error log stream status is run!")

	conf, err := confreader.LoadConfig()
	if err != nil {
		logger.Error(err.Error())
	}
	logger.Info("configs loaded")

	if err = bd.Init(conf.DMS.Host, conf.DMS.User, conf.DMS.Password, conf.DMS.DBname, conf.DMS.Port, conf.DMS.SSLmode); err != nil {
		logger.Error(err.Error())
	}

	if err = bd.Migrate(); err != nil {
		logger.Error(err.Error())
	}
	logger.Info("database is Ready ...")

	if err = hh.Init(); err != nil {
		logger.Error(err.Error())
		return
	}
	logger.Info("hh worker is OK")

	logger.Info("telegram bot worker start")
	if err := telebot.Run(conf.Tbot.API); err != nil {
		logger.Error(err.Error())
		return
	}
}
