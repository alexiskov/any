package bd

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type (
	DBentity struct {
		Socket *gorm.DB
	}

	DataEntity struct{}
)

type (
	JobAnnounce struct {
		gorm.Model
		Name           string
		Expirience     string
		Region         string
		SalaryGross    bool
		SalaryFrom     float32
		SalaryTo       float32
		SalaryCurrence string
	}
)

func (DB *DBentity) Init(host, user, password, dbname string, port int, sslmode bool) (err error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%t", host, user, password, dbname, port, sslmode)
	DB.Socket, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return
	}
	return nil
}
