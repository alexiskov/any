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

func (DB *DBentity) Init(host, user, password, dbname string, port int, sslmode bool) (err error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%t", host, user, password, dbname, port, sslmode)
	DB.Socket, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return
	}
	return nil
}
