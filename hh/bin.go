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
func SentRequest(vacancieName string, sched schedule, exp experience) (resp HHresponse, err error) {
	var hh htpcli.RequestDealer = &htpcli.HTTPclient{Socket: &http.Client{}}
	r, err := hh.NewGet("https://api.hh.ru/vacancies?text="+vacancieName, map[string]string{"User-Agent": "HH-User-Agent"}).Do()
	if err != nil {
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}

	rsp := HHresponse{}
	if err = json.Unmarshal(b, &rsp); err != nil {
		return
	}

	//resp := htpcli.HHresponse{}
	for _, v := range rsp.Items {
		if v.Experience.ID == string(exp) && v.Schedule.ID == string(sched) {
			resp.Items = append(resp.Items, v)
		}
	}

	//-----------------------------
	for _, s := range resp.Items {
		fmt.Printf("\n%+v\n", s)
	}
	//-----------------------------

	return
}
