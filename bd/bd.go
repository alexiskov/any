package bd

import (
	"errors"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

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
		City           string
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
)

var DB DBentity

func Init(host, user, password, dbname string, port int, sslmode string) (err error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s", host, user, password, dbname, port, sslmode)
	DB.Socket, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		err = fmt.Errorf("database init error: %w", err)
		return
	}
	return nil
}

func Migrate() (err error) {
	if err = DB.Socket.AutoMigrate(UserData{}, JobAnnounce{}, UserPivotVacancy{}); err != nil {
		err = fmt.Errorf("database automigration error: %w", err)
	}
	return
}

func FindOrCreateUser(tgID int64) (u UserData, err error) {
	if err = DB.Socket.Where("tg_id=?", tgID).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			u.TgID = tgID
			if err = DB.Socket.Create(&u).Error; err != nil {
				err = fmt.Errorf("user creating error: %w", err)
			}
			return
		} else {
			err = fmt.Errorf("user finding error: %w", err)
			return
		}
	}
	return
}

func (u *UserData) Update() (err error) {
	sqlu, err := FindOrCreateUser(u.TgID)
	if err != nil {
		err = fmt.Errorf("update user data error: %w", err)
		return
	}
	if u.City != "" {
		sqlu.City = u.City
	}
	if u.ExperienceYear != 0 {
		sqlu.ExperienceYear = u.ExperienceYear
	}
	if u.Schedule != 0 {
		sqlu.Schedule = u.Schedule
	}
	if u.VacancyName != "" {
		sqlu.VacancyName = u.VacancyName
	}
	if err = DB.Socket.Save(&sqlu).Error; err != nil {
		err = fmt.Errorf("updating userData error: %w", err)
		return
	}
	return nil
}
