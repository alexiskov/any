package main

import (
	"fmt"
	"os"

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
	logger.Info("\nconfigs loaded\n")

	go tgBot(conf.Tbot.API)
	logger.Info("telegram bot-worker is run...")

	res, err := hh.SentRequest("golang", hh.REMOTE_JOB, hh.NO_EXPERIENCE, 0)
	if err != nil {
		logger.Error(err.Error())
	}

	for _, vac := range res.Items {
		fmt.Printf("\n%+v\n", vac)
	}

	for {
	}
}

func tgBot(tgAPIkey string) {
	if err := telebot.Run(tgAPIkey); err != nil {
		logger.Error(err.Error())
	}
}
