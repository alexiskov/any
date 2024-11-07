package hh

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"
	"vacancydealer/bd"
	"vacancydealer/htpcli"
	"vacancydealer/logger"
)

func ConvertSerchPatternModelDBtoHH(from bd.VacancyNamePatterns) (to []HHfilterData) {
	for _, v := range from {
		to = append(to, HHfilterData{VacancyName: v.VacancyName})
	}
	return
}

func (hh HHresponse) ConvertItemsToDB() (bdja bd.JobAnnounces) {
	for _, vac := range hh.Items {
		id, err := strconv.Atoi(vac.ID)
		if err != nil {
			continue
		}

		locID, err := strconv.Atoi(vac.Area.RegionID)
		if err != nil {
			continue
		}

		bdja = append(bdja, bd.JobAnnounce{ItemId: id, Name: vac.Name, Company: vac.Employer.Name, Area: locID, Expierence: vac.Experience.Name, SalaryGross: vac.Salary.Gross, SalaryFrom: vac.Salary.From, SalaryTo: vac.Salary.To, SalaryCurrency: vac.Salary.Currency, PublishedAt: vac.PublishedAt, Requirement: vac.Snippet.Requirement, Responsebility: vac.Snippet.Responsibility, Link: vac.PageURL})
	}
	return
}

// --------------------------------------------------------------------------------------------------- ProdMethod method to hhAPI query due
func (hf HHfilterData) GetJobAnnounces() (resp HHresponse, err error) {
	var hh htpcli.RequestDealer = &htpcli.HTTPclient{Socket: &http.Client{}}
	urq := "https://api.hh.ru/vacancies?applicant_comments_order=creation_time_desc&per_page=100"

	if hf.VacancyName != "" {
		urq += "&text=NAME%3A(" + hf.VacancyName + ")"
	}

	r, err := hh.NewGet(urq, map[string]string{"User-Agent": "HH-User-Agent"}).Do()
	if err != nil {
		return
	}

	switch r.StatusCode {
	case http.StatusBadRequest:
		err = StatusBadRequest
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}

	if err = json.Unmarshal(b, &resp); err != nil {
		return
	}
	return
}

// vacancy announce query to HHunter-API send
func WorkerStart(pauseDuration int) {
	for {
		keys, err := bd.GetVacancyPatterns()
		if err != nil {
			logger.Error(err.Error())
			continue
		}
		for _, k := range ConvertSerchPatternModelDBtoHH(keys) {
			resp, err := k.GetJobAnnounces()
			if err != nil {
				logger.Error(err.Error())
				continue
			}
			if err = resp.ConvertItemsToDB().SaveInDB(); err != nil {
				logger.Error(err.Error())
				continue
			}
			time.Sleep(time.Duration(pauseDuration/len(keys)) * time.Second)
		}

		time.Sleep(time.Duration(pauseDuration) * time.Second)
	}

}
