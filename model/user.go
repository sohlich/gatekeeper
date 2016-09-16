package model

import (
	"fmt"
	"log"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model `valid:"-"`
	UserID     string `json:"userid"`
	Email      string `valid:"email"`
	Password   string `valid:"-"`
	Firstname  string `json:"firstname"`
	Lastname   string `json:"lastname"`
	Activated  bool   `valid:"-"`
	Expiration int64
	LastAccess int64
}

type UserManager interface {
	Save(u *User) error
	FindByID(id uint) (*User, error)
	FindByUserID(userID string) (*User, error)
}

type SqlUserManager struct {
	db *gorm.DB
}

func (s *SqlUserManager) Save(u *User) error {
	return s.db.Save(u).Error
}

func (s *SqlUserManager) FindByID(id uint) (u *User, err error) {
	u = &User{}
	err = s.db.First(u, id).Error
	return
}

func (s *SqlUserManager) FindByUserID(userID string) (u *User, err error) {
	u = &User{}
	err = s.db.Where("userid = ?", userID).First(u).Error
	return
}

func (u *User) Stringer() string {
	if u != nil {
		return u.UserID
	}
	return ""
}

func (u *User) TableName() string {
	return "user"
}

func (u *User) CleaUp() {
	u.Password = ""
}

func CreateUser(u *User) (err error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.Wrap(err, "Cannot create user")
	}
	valid, err := govalidator.ValidateStruct(u)
	if err != nil {
		return errors.Wrap(err, "Cannot validate user")
	}
	if !valid {
		return fmt.Errorf("Object data are missing")
	}

	u.Expiration = time.Now().Add(time.Hour * 8765).Unix()
	u.Activated = false

	log.Printf("Got password: %v\n", bytes)
	u.Password = string(bytes)
	err = userStorage.Save(u)
	return errors.Wrap(err, "Cannot create user")
}

func FindUserByUserid(userID string) (user User) {
	db.Where("userid = ?", userID).First(&user)
	return
}

func FindUserByID(id uint) (*User, error) {
	return userStorage.FindByID(id)
}

func SaveUser(u *User) error {
	return userStorage.Save(u)
}

func FindUserByActivity(a *Activity) (*User, error) {
	return userStorage.FindByID(a.UserID)
}
