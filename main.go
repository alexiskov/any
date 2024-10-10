package main

import (
	"os"

	"vacancydealer/logger"
)

func main() {
	logger.InitInfoTextlog(os.Stdout)
	logger.Info("logger status is Run...")

	logger.InitDebugJSONlog(os.Stdout)
	logger.Info("debug log stream status is run!")

}
