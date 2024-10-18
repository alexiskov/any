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
)
