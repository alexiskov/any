package main

import (
	"fmt"
	"os"
	"project/confreader"
	"project/htpsrv"
	"project/logger"
	"strconv"
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
	logger.Info("configs loaded.")

	logger.Info("webservice is ON, port: " + strconv.Itoa(configs.WebServer.Port))
	htpsrv.Start(configs.WebServer.Port)

	fmt.Printf("%+v\n", configs.DMS)
}
