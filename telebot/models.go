package telebot

import "time"

type (
	UserData struct {
		TgID            int64
		Vacancy         string
		Location        string
		Schedule        string
		ExperienceYears int
	}

	UserStateData struct {
		State uint8
		User  UserData
		Date  time.Time
	}

	ScheduleType struct {
		Name  string
		Valie int
	}

	JobAnnounce struct {
		Name           string
		Company        string
		Area           string
		Expierence     string
		SalaryGross    bool
		SalaryFrom     float64
		SalaryTo       float64
		SalaryCurrency string
		PublishedAt    string
		Schedule       string
		Requirement    string
		Responsebility string
		Link           string
	}
)
