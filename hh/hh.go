package hh

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"vacancydealer/bd"
	"vacancydealer/htpcli"
)

type (
	experience string
	schedule   string
)

var (
	StatusBadRequest = errors.New("status BadRequest")
)

// Инициализация базовых справочников из ХэХа
// Получение данных Локаций
// Получение графиков работ
// Запись в БД
func Init() (err error) {
	areasHH, err := getAreas()
	if err != nil {
		return
	}

	if err = areasHH.CreateToDB(); err != nil {
		return
	}

	schedulesHH, err := GetSchedulesList()
	if err != nil {
		return
	}
	if err = schedulesHH.SchedulesModelConvert().CreateToDB(); err != nil {
		return
	}

	return nil
}

// sent query to HH
func (dataFilter UserFilter) GetVacancies(pp, page int) (rsp HHresponse, err error) {
	var hh htpcli.RequestDealer = &htpcli.HTTPclient{Socket: &http.Client{}}
	urq := fmt.Sprintf("https://api.hh.ru/vacancies?&experience=%s&schedule=%s&applicant_comments_order=creation_time_desc&per_page=%d", dataFilter.Experience, dataFilter.Schedule, pp)
	if dataFilter.Vacancyname != "" {
		if strings.Contains(dataFilter.Vacancyname, " ") {
			dataFilter.Vacancyname = strings.ReplaceAll(dataFilter.Vacancyname, " ", "+")
		}
		urq += "&text=NAME%3A(" + dataFilter.Vacancyname + ")"
	}

	if page != 0 {
		urq += "&page=" + strconv.Itoa(page)
	}
	if dataFilter.Location != 0 {
		urq += "&area=" + strconv.Itoa(dataFilter.Location)
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

	if err = json.Unmarshal(b, &rsp); err != nil {
		return
	}
	return
}

// query to HH API
// Получаем локации от ХэХа
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

// query to HH API
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

// Справочник локаций из ХэХа
// Обработка и запись в БД
//
//	//..Обработка стран
//
// *Принятая рессивером с ХэХа схема json разбирается циклом
// Разведение стран, областей и городов по разным справочникам
// ..Обработка локаций-
func (areasHH Countries) CreateToDB() (err error) {
	sqlcountries := bd.Countries{}
	sqlregions := bd.Regions{}
	sqlcities := bd.Cities{}

	// [htym]!
	type Shit struct {
	}

	for _, country := range areasHH {
		coi, err := strconv.Atoi(country.ID)
		if err != nil {
			err = fmt.Errorf("regions on DB create, region id parse error: %w", err)
			return err
		}
		sqlcountries = append(sqlcountries, bd.Country{ID: uint(coi), Name: country.Name})

		for _, region := range country.AreaList {
			ri, err := strconv.Atoi(region.ID)
			if err != nil {
				err = fmt.Errorf("regions on DB create, region id parse error: %w", err)
				return err
			}

			////Обработка городов--
			//.отсев регионов не содержащих городов
			//.Отсеятся ,,МЕгаполисы???(не имеют родителя области. Имеют страну))))
			if len(region.AreaList) != 0 { //Отбираем регионы не содержащие городов
				sqlregions = append(sqlregions, bd.Region{ID: uint(ri), Name: region.Name, Owner: uint(coi)}) //замена
				for _, city := range region.AreaList {
					ciID, err := strconv.Atoi(city.ID)
					if err != nil {
						err = fmt.Errorf("regions on DB create, region id parse error: %w", err)
						return err
					}
					sqlcities = append(sqlcities, bd.City{ID: uint(ciID), Name: city.Name, Owner: uint(ri)})
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

// package HH model to model of DB package convert
func (from ScheduleData) SchedulesModelConvert() (to bd.Schedules) {
	for _, s := range from.List {
		to = append(to, bd.Schedule{HhID: s.Id, Name: s.Name})
	}
	return
}

// user convert model data of users from package bd to models of UserFilter
func ConvertUserData(userdata []bd.UserData) (userFilterList []UserFilter) {
	for _, bdUd := range userdata {
		userFilterTemp := UserFilter{TgID: bdUd.TgID, Vacancyname: bdUd.VacancyName, Location: int(bdUd.Location), Schedule: bdUd.Schedule}
		if bdUd.ExperienceYear < 1 {
			userFilterTemp.Experience = "noExperience"
		} else if bdUd.ExperienceYear > 0 && bdUd.ExperienceYear < 4 {
			userFilterTemp.Experience = "between1And3"
		} else if bdUd.ExperienceYear > 3 && bdUd.ExperienceYear < 7 {
			userFilterTemp.Experience = "between3And6"
		} else if bdUd.ExperienceYear > 6 {
			userFilterTemp.Experience = "moreThan6"
		}
		userFilterList = append(userFilterList, userFilterTemp)
	}
	return
}
