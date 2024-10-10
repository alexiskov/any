package execute

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"project/htpcli"
)

func StartHH() error {
	cli := htpcli.HTTPclient{Socket: &http.Client{}}
	resp, err := cli.NewGet("https://api.hh.ru/vacancies?text=golang", map[string]string{"User-Agent": "HH-User-Agent"}).Do()
	if err != nil {
		return err
	}

	if b, err := io.ReadAll(resp.Body); err != nil {
		return err
	} else {
		rsp := htpcli.HHresponse{}
		if err = json.Unmarshal(b, &rsp); err != nil {
			return err
		}
		for _, v := range rsp.Items {
			if v.Experience.ID == "noExperience" && v.Schedule.ID == "remote" {
				fmt.Printf("%+v\n\n\n", v)
			}
		}
	}
	return nil
}
