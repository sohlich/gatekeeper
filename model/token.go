package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Token struct {
	gorm.Model
	UserID     uint
	RefToken   string
	JwtToken   string
	Expiration int64
}

func (t *Token) TableName() string {
	return "token"
}

type TokenManager interface {
	Save(t *Token) error
	FindByRef(refTkn string) (*Token, error)
	FindActiveByUser(u *User) (*Token, error)
	InvalidateAllByUserID(id uint) error
}

type SQLTokenManager struct {
	db *gorm.DB
}

func (s *SQLTokenManager) Save(t *Token) error {
	return s.db.Save(t).Error
}

func (s *SQLTokenManager) FindByRef(refTkn string) (*Token, error) {
	t := &Token{}
	err := s.db.First(t, "ref_token = ?", refTkn).Error
	return t, err
}

func (s *SQLTokenManager) FindActiveByUser(u *User) (*Token, error) {
	t := &Token{}
	err := s.db.First(t, "user_id = ? and expiration >= ? ", u.ID, time.Now().Unix()).Error
	return t, err
}

func (s *SQLTokenManager) InvalidateAllByUserID(userID uint) error {
	now := time.Now().Unix()
	return s.db.Exec("update token set expiration = ? where user_id = ? and expiration >= ?", now, userID, now).Error
}

func SaveToken(t *Token) error {
	return tokenStorage.Save(t)
}

func FindActiveByUser(u *User) (*Token, error) {
	return tokenStorage.FindActiveByUser(u)
}

func FindTokenByRef(refTkn string) (*Token, error) {
	return tokenStorage.FindByRef(refTkn)
}

func InvalidateAllSessionByUsedID(ID uint) error {
	return tokenStorage.InvalidateAllByUserID(ID)
}
