package model

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
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

func NewToken(user User, jwtToken string) Token {
	token := Token{}
	token.JwtToken = jwtToken
	token.UserID = user.ID
	token.RefToken = uuid.NewV4().String()
	token.Expiration = time.Now().Add(time.Hour * 72).Unix()
	return token
}
