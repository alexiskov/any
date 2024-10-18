package bd

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB DBentity

// ---------------------------------------->>>INITIALIZATION---------------------------------------------------------------------
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
	if err = DB.Socket.AutoMigrate(UserData{}, JobAnnounce{}, UserPivotVacancy{}, Country{}, Region{}, City{}); err != nil {
		err = fmt.Errorf("database automigration error: %w", err)
	}
	return
}

// ----------------------------------------<<<INITIALIZATION----------------------------------------------------------------------

// ------------------------------------------------------------->>>LOCATION WRITERS-----------------------------------------------------
func (countries Countries) WriteCountries() (err error) {
	if err = DB.Socket.Save(&countries).Error; err != nil {
		err = fmt.Errorf("list of region save error: %w", err)
	}
	return
}

func (regions Regions) WriteRegions() (err error) {
	if err = DB.Socket.Save(&regions).Error; err != nil {
		err = fmt.Errorf("list of region save error: %w", err)
	}
	return
}

func (cities Cities) WriteCities() (err error) {
	if err = DB.Socket.Save(&cities).Error; err != nil {
		err = fmt.Errorf("list of region save error: %w", err)
	}
	return
}

// -------------------------------------------------------------<<<LOCATION WRITERS-----------------------------------------------------

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

// User Data Update
func (u UserData) Update() (err error) {
	sqlu, err := FindOrCreateUser(u.TgID)
	if err != nil {
		err = fmt.Errorf("update user data error: %w", err)
		return
	}
	sqlu.ExperienceYear = u.ExperienceYear
	sqlu.Schedule = u.Schedule
	sqlu.VacancyName = u.VacancyName
	if err = DB.Socket.Model(&sqlu).Select("schedule", "vacancy_name", "experience_year").Updates(UserData{Schedule: sqlu.Schedule, VacancyName: sqlu.VacancyName, ExperienceYear: sqlu.ExperienceYear}).Error; err != nil {
		err = fmt.Errorf("updating userData error: %w", err)
		return
	}
	return nil
}

func (u UserData) UpdateLocation() (err error) {
	if err = DB.Socket.Model(&u).Where("tg_id=?", u.TgID).Update("location", u.Location).Error; err != nil {
		err = fmt.Errorf("user data location on db update error: %w", err)
	}
	return
}

func FindCitiesByName(cityName string) (cities Cities, err error) {
	if err = DB.Socket.Where("LOWER(name) like ?", "%"+strings.ToLower(cityName)+"%").Find(&cities).Error; err != nil {
		err = fmt.Errorf("cities by name finding error: %w", err)
		return
	}
	IDs := make([]uint, 0, len(cities))
	for _, city := range cities {
		IDs = append(IDs, city.Owner)
	}

	regions := Regions{}
	if err = DB.Socket.Where("id in ?", IDs).Find(&regions).Error; err != nil {
		err = fmt.Errorf("regions by id finding error: %w", err)
		return
	}
	IDs = make([]uint, 0, len(regions))
	for _, region := range regions {
		IDs = append(IDs, region.Owner)
	}

	countries := Countries{}
	if err = DB.Socket.Where("id in ?", IDs).Find(&countries).Error; err != nil {
		err = fmt.Errorf("countries by id finding error: %w", err)
		return
	}

	for i, city := range cities {
		for _, region := range regions {
			if city.Owner == region.ID {
				cities[i].Name = region.Name + ", " + city.Name
				for _, country := range countries {
					if region.Owner == country.ID {
						cities[i].Name = country.Name + ", " + cities[i].Name
					}
				}
			}
		}
	}
	return
}

func FindRegionByName(regionName string) (regions Regions, err error) {
	if err = DB.Socket.Where("LOWER(name) like ?", "%"+strings.ToLower(regionName)+"%").Find(&regions).Error; err != nil {
		err = fmt.Errorf("Find region by name error: %w", err)
		return
	}
	IDs := make([]uint, 0, len(regions))
	for _, region := range regions {
		IDs = append(IDs, region.Owner)
	}

	countries := Countries{}
	if err = DB.Socket.Where("id in ?", IDs).Find(&countries).Error; err != nil {
		err = fmt.Errorf("find countries by id error: %w", err)
		return
	}

	for i, region := range regions {
		for _, country := range countries {
			if region.Owner == country.ID {
				regions[i].Name = country.Name + ", " + region.Name
			}
		}
	}
	return
}

func FindCountries() (countries Countries, err error) {
	if err = DB.Socket.Find(&countries).Error; err != nil {
		err = fmt.Errorf("countries finding error: %w", err)
	}
	return
}

// Поиск локации по ИД
// Проверяет ИД по порядку в таблицах: стран, регионов, населенных пунктов
func FindLocByID(locID uint) (locName string, err error) {
	locName = "не имеет значения"
	country := Country{}
	if err = DB.Socket.Where("id=?", locID).First(&country).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			region := Region{}
			if err = DB.Socket.Where("id=?", locID).First(&region).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					city := City{}
					if err = DB.Socket.Where("id=?", locID).First(&city).Error; err != nil {
						if !errors.Is(err, gorm.ErrRecordNotFound) {
							err = fmt.Errorf("location by ID finding error: %w", err)
							return
						}
					} else {
						locName = city.Name
					}
				} else {
					err = fmt.Errorf("location by ID finding error: %w", err)
					return
				}

			} else {
				locName = region.Name
			}

		} else {
			err = fmt.Errorf("location by ID finding error: %w", err)
			return
		}

	} else {
		locName = country.Name
	}
	return
}
