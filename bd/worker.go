package bd

import (
	"fmt"
	"strings"
	"vacancydealer/logger"
)

var (
	WorkDue = make(chan bool)
)

func GetAllUserData() (ud UserDataList, err error) {
	if err = DB.Socket.Find(&ud).Error; err != nil {
		err = fmt.Errorf("al user data getting error: %w", err)
	}
	return
}

func (ud UserDataList) MakeVacNameSearchPatternPOOL() (searchKeys VacancyNamePatterns) {
	tempNameList := make(map[string]bool, len(ud))
	for _, d := range ud {
		if _, ok := tempNameList[d.VacancyName]; !ok {
			tempNameList[d.VacancyName] = true
		} else {
			for k, v := range tempNameList {
				if v && strings.Replace(strings.ToLower(k), " ", "", -1) == strings.Replace(strings.ToLower(d.VacancyName), " ", "", -1) {
					continue
				}
			}
		}

		for subk, v := range tempNameList {
			if v {
				for k, v1 := range tempNameList {
					if v1 && strings.Contains(strings.ToLower(k), strings.ToLower(subk)) && subk != "" {
						delete(tempNameList, k)
						tempNameList[subk] = true
					}
				}
			}
		}
	}

	var sqlVacID uint = 1
	for k, v := range tempNameList {
		if v {
			searchKeys = append(searchKeys, VacancynameSearchPattern{ID: sqlVacID, VacancyName: k})
			sqlVacID++
		}
	}

	return
}

func (searchKeys VacancyNamePatterns) SaveInDB() (err error) {
	if err = DB.Socket.Save(&searchKeys).Error; err != nil {
		err = fmt.Errorf("vacancy name pool in database saving error: %w", err)
	}
	return
}

func StarWorker(ch <-chan bool) {
	for {
		select {
		case <-ch:
			ud, err := GetAllUserData()
			if err != nil {
				logger.Error(err.Error())
				continue
			}
			if err = ud.MakeVacNameSearchPatternPOOL().SaveInDB(); err != nil {
				logger.Error(err.Error())
				continue
			}
		}
	}
}
