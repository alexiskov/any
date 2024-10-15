package main

import (
	"fmt"
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

	res, err := hh.SentRequest("golang", hh.REMOTE_JOB, hh.NO_EXPERIENCE, 0)
	if err != nil {
		logger.Error(err.Error())
	}

	for _, vac := range res.Items {
		fmt.Printf("\n%+v\n", vac)
	}

	if err := telebot.Run(conf.Tbot.API); err != nil {
		logger.Error(err.Error())
		return
	}
}
