package hh

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
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
func SentRequest(vacancieName string, sched schedule, exp experience, page int) (rsp HHresponse, err error) {
	var hh htpcli.RequestDealer = &htpcli.HTTPclient{Socket: &http.Client{}}
	urq := fmt.Sprintf("https://api.hh.ru/vacancies?text=%s&experience=%s&schedule=%s&applicant_comments_order=creation_time_desc&per_page=100", vacancieName, exp, sched)
	if page != 0 {
		urq += "&page=" + strconv.Itoa(page)
	}
	r, err := hh.NewGet(urq, map[string]string{"User-Agent": "HH-User-Agent"}).Do()
	if err != nil {
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}

	if err = json.Unmarshal(b, &rsp); err != nil {
		return
	}
	return
}
