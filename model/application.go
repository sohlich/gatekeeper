package model

import "github.com/jinzhu/gorm"

type Application struct {
	gorm.Model
	Name   string
	APIKey string
	Owner  string
}

func (a *Application) TableName() string {
	return ""
}
