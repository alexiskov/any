package bd_test

import (
	"fmt"
	"testing"
	"vacancydealer/bd"
)

func TestMakeVacNameSearchPatternPOOL(t *testing.T) {
	moka := bd.UserDataList{
		{VacancyName: "Golang"},
		{VacancyName: "go"},
		{VacancyName: "массажист лиц"},
		{VacancyName: "массаж"},
		{VacancyName: "Golang developer"},
		{VacancyName: "Develop"},
		{VacancyName: "backend developer"},
		{VacancyName: "back"},
		{VacancyName: ""},
		{VacancyName: "lo"},
	}

	if len(moka.MakeVacNameSearchPatternPOOL()) != 5 {
		t.Errorf("Result was incorrect, expected %d, got %d", 5, len(moka.MakeVacNameSearchPatternPOOL()))
	}

	fmt.Println(moka.MakeVacNameSearchPatternPOOL())

	for _, v := range moka.MakeVacNameSearchPatternPOOL() {
		switch v.VacancyName {
		case "lo":
			continue
		case "":
			continue
		case "go":
			continue
		case "back":
			continue
		case "массаж":
			continue
		}

		t.Errorf("Result was incorrect: %s", v.VacancyName)
	}

}
