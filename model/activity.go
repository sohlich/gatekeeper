package model

import "github.com/jinzhu/gorm"

const (
	PasswordReset  = "PASSWORD_RESET"
	UserActivation = "USER_ACTIVATION"
)

type Activity struct {
	gorm.Model
	Type       string
	Token      string
	UserID     uint
	Time       int64
	Expiration int64
	Used       int64
}

func (a *Activity) TableName() string {
	return "activity"
}

func FindActivityByToken(tkn string) (*Activity, error) {
	act := &Activity{}
	err := db.Where("token = ?", tkn).First(act).Error
	return act, err
}

func CreateActivity(a *Activity) error {
	return db.Begin().Create(a).
		Commit().Error
}

func UpdateActivity(a *Activity) error {
	return db.Save(a).Error
}
