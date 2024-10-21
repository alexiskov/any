package hh

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"vacancydealer/bd"
	"vacancydealer/htpcli"
)

type (
	experience string
	schedule   string
)

// sent query to HH
func GetVacancy(name string, sched schedule, exp experience, page int) (rsp HHresponse, err error) {
	var hh htpcli.RequestDealer = &htpcli.HTTPclient{Socket: &http.Client{}}
	urq := fmt.Sprintf("https://api.hh.ru/vacancies?text=%s&experience=%s&schedule=%s&applicant_comments_order=creation_time_desc&per_page=100", name, exp, sched)
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

func Init() (err error) {
	r, err := getAreas()
	if err != nil {
		return
	}

	if err = r.CreateToDB(); err != nil {
		return
	}

	sch, err := GetSchedulesList()
	if err != nil {
		return
	}
	if err = sch.SchedulesModelConvert().CreateToDB(); err != nil {
		return
	}

	return nil
}

func getAreas() (rsp Countries, err error) {
	var hh htpcli.RequestDealer = &htpcli.HTTPclient{Socket: &http.Client{}}
	urq := "https://api.hh.ru/areas"
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

func (countries Countries) CreateToDB() (err error) {
	sqlcountries := bd.Countries{}
	sqlregions := bd.Regions{}
	sqlcities := bd.Cities{}

	re, err := regexp.Compile(`\(.*\)`)
	if err != nil {
		err = fmt.Errorf("Create logations on DB -> regxp pattern compilation error: %w", err)
		return err
	}

	for _, country := range countries {
		coi, err := strconv.Atoi(country.ID)
		if err != nil {
			err = fmt.Errorf("regions on DB create, region id parse error: %w", err)
			return err
		}
		sqlcountries = append(sqlcountries, bd.Country{ID: uint(coi), Name: country.Name})

		for _, region := range country.Regions {
			ri, err := strconv.Atoi(region.ID)
			if err != nil {
				err = fmt.Errorf("regions on DB create, region id parse error: %w", err)
				return err
			}
			rgxRegion := re.ReplaceAllString(region.Name, "")
			if len(region.Cities) != 0 {
				sqlregions = append(sqlregions, bd.Region{ID: uint(ri), Name: rgxRegion, Owner: uint(coi)})
				for _, city := range region.Cities {
					ci, err := strconv.Atoi(city.ID)
					if err != nil {
						err = fmt.Errorf("regions on DB create, region id parse error: %w", err)
						return err
					}

					rgxCity := re.ReplaceAllString(city.Name, "")
					sqlcities = append(sqlcities, bd.City{ID: uint(ci), Name: rgxCity, Owner: uint(ri)})
				}
			} else {
				sqlcities = append(sqlcities, bd.City{ID: uint(ri), Name: region.Name, Owner: uint(coi)})
			}

		}
	}

	if err = sqlcountries.WriteCountries(); err != nil {
		return
	}
	if err = sqlregions.WriteRegions(); err != nil {
		return
	}
	if err = sqlcities.WriteCities(); err != nil {
		return
	}

	return nil
}

func GetSchedulesList() (rsp ScheduleData, err error) {
	var hh htpcli.RequestDealer = &htpcli.HTTPclient{Socket: &http.Client{}}
	urq := "https://api.hh.ru/dictionaries"
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

// package HH model to model of DB package convert
func (from ScheduleData) SchedulesModelConvert() (to bd.Schedules) {
	for _, s := range from.List {
		to = append(to, bd.Schedule{HhID: s.Id, Name: s.Name})
	}
	return
}
