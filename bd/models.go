package bd

import "gorm.io/gorm"

type (
	DBentity struct {
		Socket *gorm.DB
	}

	DataEntity struct{}
)

type (
	UserData struct {
		gorm.Model
		TgID           int64 `gorm:"uniqueIndex"`
		VacancyName    string
		ExperienceYear int
		Schedule       string
		Location       uint
	}

	UserDataList []UserData

	JobAnnounce struct {
		ItemId         int    `gorm:"primaryKey"`
		Name           string `gorm:"index"`
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

	JobAnnounces []JobAnnounce

	UserPivotVacancy struct {
		gorm.Model
		UID    uint
		JobID  uint
		Showed bool
	}

	Country struct {
		ID   uint   `gorm:"primaryKey"`
		Name string `gorm:"index"`
	}
	Region struct {
		ID    uint   `gorm:"primaryKey"`
		Name  string `gorm:"index"`
		Owner uint
	}
	City struct {
		ID    uint   `gorm:"primaryKey"`
		Name  string `gorm:"index"`
		Owner uint
	}
	Countries []Country
	Regions   []Region
	Cities    []City

	Schedule struct {
		HhID string `gorm:"primaryKey"`
		Name string
	}

	Schedules []Schedule

	VacancynameSearchPattern struct {
		VacancyName string `gorm:"index"`
	}
)
