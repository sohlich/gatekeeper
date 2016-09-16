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

type ActivityManager interface {
	Save(a *Activity) error
	FindByToken(tkn string) (*Activity, error)
}

type SQLActivityManager struct {
	db *gorm.DB
}

func (s *SQLActivityManager) Save(a *Activity) error {
	return db.Save(a).Error
}

func (s *SQLActivityManager) FindByToken(tkn string) (a *Activity, err error) {
	a = &Activity{}
	err = db.Where("token = ?", tkn).First(a).Error
	return
}

func FindActivityByToken(tkn string) (*Activity, error) {
	return activityStorage.FindByToken(tkn)
}

func CreateActivity(a *Activity) error {
	return activityStorage.Save(a)
}

func UpdateActivity(a *Activity) error {
	return activityStorage.Save(a)
}
