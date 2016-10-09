package model

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

var db *gorm.DB

var userStorage UserManager
var activityStorage ActivityManager
var tokenStorage TokenManager

type Table interface {
	TableName() string
}

func InitDB() (err error) {
	db, err = gorm.Open("mysql",
		"user:user@/attendence?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		return errors.Wrap(err, "Cannot connect to database")
	}
	log.Println("Doing health check...")
	tables := []Table{
		&User{},
		&Activity{},
		&Token{}}

	for _, val := range tables {
		if !db.HasTable(val) {
			db.AutoMigrate(val)
		}
	}

	userStorage = &SqlUserManager{
		db,
	}
	activityStorage = &SQLActivityManager{
		db,
	}
	tokenStorage = &SQLTokenManager{
		db,
	}

	log.Println("Database tables OK")
	return nil
}

func CloseDB() error {
	return db.Close()
}
