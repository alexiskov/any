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
		Schedule       int
		Location       uint
	}
	JobAnnounce struct {
		gorm.Model
		ItemId         int `gorm:"uniqueIndex"`
		Name           string
		Expierence     string
		SalaryGross    bool
		SalaryFrom     int
		SalaryTo       int
		SalaryCurrency string
		PublishedAt    int64
		Schedule       string
		Requirement    string
		Responsebility string
		Link           string
	}
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
)
