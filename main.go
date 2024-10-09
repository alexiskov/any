package main

import (
	"fmt"
	"os"
	"project/confreader"
	"project/logger"
)

func main() {
	logger.InitInfoTextlog(os.Stdout)
	logger.Info("logger status is Run...")

	logger.InitDebugJSONlog(os.Stdout)
	logger.Info("debug log stream status is run!")

	configs, err := confreader.LoadConfig()
	if err != nil {
		logger.Debug(err.Error())
		return
	}

	fmt.Printf("%+v\n", configs.DMS)
}
