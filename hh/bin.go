package hh

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"vacancydealer/htpcli"
)

const (
	//Experience values HH
	NO_EXPERIENCE   experience = "noExperience"
	BETWEN_1_3_YEAR experience = "between1And3"
	BETWEN_3_6_YEAR experience = "between3And6"

	//Schedule values HH
	REMOTE_JOB schedule = "remote"
)

type (
	experience string
	schedule   string
)

// sent query to HH
func Start(vacancieName string, sched schedule, exp experience) error {
	hh := htpcli.HTTPclient{Socket: &http.Client{}}
	resp, err := hh.NewGet("https://api.hh.ru/vacancies?text=golang&period=1", map[string]string{"User-Agent": "HH-User-Agent"}).Do()
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
			if v.Experience.ID == string(exp) && v.Schedule.ID == string(sched) {
				fmt.Printf("%+v\n\n\n", v)
			}
		}
	}
	return nil
}
