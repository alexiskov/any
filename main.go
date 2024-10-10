package main

import (
	"os"
	"project/execute"
	"project/logger"
)

func main() {
	logger.InitInfoTextlog(os.Stdout)
	logger.Info("logger status is Run...")

	logger.InitDebugJSONlog(os.Stdout)
	logger.Info("debug log stream status is run!")

	if err := execute.StartHH(); err != nil {
		logger.Debug(err.Error())
	}
}
