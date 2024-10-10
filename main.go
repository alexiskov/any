package main

import (
	"os"

	"vacancydealer/hh"
	"vacancydealer/logger"
)

func main() {
	logger.InitInfoTextlog(os.Stdout)
	logger.Info("logger status is Run...")

	logger.InitDebugJSONlog(os.Stdout)
	logger.Info("debug log stream status is run!")

	if err := hh.Start("golang", hh.REMOTE_JOB, hh.NO_EXPERIENCE); err != nil {
		logger.Debug(err.Error())
	}
}
