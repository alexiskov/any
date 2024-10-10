package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"project/htpcli"
	"project/logger"
)

func main() {
	logger.InitInfoTextlog(os.Stdout)
	logger.Info("logger status is Run...")

	logger.InitDebugJSONlog(os.Stdout)
	logger.Info("debug log stream status is run!")

	cli := htpcli.HTTPclient{Socket: &http.Client{}}
	resp, err := cli.NewGet("https://api.hh.ru/vacancies", map[string]string{"User-Agent": "HH-User-Agent"}).Do()
	if err != nil {
		logger.Debug(err.Error())
	}

	if b, err := io.ReadAll(resp.Body); err != nil {
		logger.Debug(err.Error())
	} else {
		fmt.Printf("%s", string(b))
	}

}
