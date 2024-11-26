package hh

import (
	"encoding/json"
	"fmt"
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

		bd.FindLocByID(uint(locID))

		bdja = append(bdja, bd.JobAnnounce{ItemId: uint(id), Name: vac.Name, Company: vac.Employer.Name, Area: int(city.ID), Region: int(region.ID), Country: int(country.ID), Expierence: vac.Experience.ID, SalaryGross: vac.Salary.Gross, SalaryFrom: vac.Salary.From, SalaryTo: vac.Salary.To, SalaryCurrency: vac.Salary.Currency, PublishedAt: vac.PublishedAt, Schedule: vac.Schedule.ID, Requirement: vac.Snippet.Requirement, Responsebility: vac.Snippet.Responsibility, Link: vac.PageURL})
	}
	return
}

func Reader(r *http.Response) (dataBytes []byte, err error) {
	switch r.StatusCode {
	case http.StatusBadRequest:
		err := r.Body.Close()
		return nil, fmt.Errorf("%d bad request err: %w", http.StatusBadRequest, err)
	case http.StatusOK:
		dataBytes, err = io.ReadAll(r.Body)
		if err != nil {
			dataBytes = nil
		}
		return
	}
	return dataBytes, nil
}

// --------------------------------------------------------------------------------------------------- ProdMethod method to hhAPI query due
func (hf HHfilterData) GetJobAnnounces() (hhResponseRest HHresponse, err error) {
	if hf.VacancyName != "" {
		uRqPreset := "https://api.hh.ru/vacancies?applicant_comments_order=creation_time_desc&per_page=100"

		uRqPreset += "&text=NAME%3A(" + hf.VacancyName + ")"
		var hh htpcli.RequestDealer = &htpcli.HTTPclient{Socket: &http.Client{}}
		getResp, err := hh.NewGet(uRqPreset, map[string]string{"User-Agent": "HH-User-Agent"}).Do()
		if err != nil {
			return hhResponseRest, err
		}

		d, err := Reader(getResp)
		json.Unmarshal(d, &hhResponseRest)

		return hhResponseRest, nil
	}

	return
}

// vacancy announce query to HHunter-API send
func WorkerStart(pauseDuration int) {
	time.Sleep(time.Duration(10) * time.Second)

	areas, err := bd.CountriesList()
	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}

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
			if err = resp.ConvertItemsToDB(areas).SaveInDB(); err != nil {
				logger.Error(err.Error())
				continue
			}
			if len(keys) != 0 {
				time.Sleep(time.Duration(pauseDuration/len(keys)) * time.Second)
			}

		}

		time.Sleep(time.Duration(pauseDuration) * time.Second)
	}

}
