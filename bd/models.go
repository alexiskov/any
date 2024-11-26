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
		ItemId         uint   `gorm:"primaryKey"`
		Name           string `gorm:"index"`
		Company        string
		Area           int
		Region         int
		Country        int
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
		UID   uint
		JobID uint `gorm:"uniqueIndex"`
	}

	CountrySQL struct {
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

	AreaData struct {
		Countries []CountrieModel
	}
	CountrieModel struct {
		Count   AreaEntity
		Regions Regions
	}
	RegionModel struct {
		Region AreaEntity
		Cities Cities
	}

	AreaEntity struct {
		ID    uint
		Name  string
		Owner uint
	}

	Countries []CountrieModel
	Regions   []RegionModel
	Cities    []AreaEntity

	Schedule struct {
		HhID string `gorm:"primaryKey"`
		Name string
	}

	Schedules []Schedule

	VacancynameSearchPattern struct {
		ID          uint   `gorm:"primaryKey"`
		VacancyName string `gorm:"index"`
	}

	VacancyNamePatterns []VacancynameSearchPattern
)
