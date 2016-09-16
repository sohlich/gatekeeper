package logic

import (
	"time"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/sohlich/gatekeeper/mail"
	"github.com/sohlich/gatekeeper/model"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(user *model.User) error {
	if err := model.CreateUser(user); err != nil {
		return err
	}
	user.CleaUp()
	now := time.Now()
	activity := &model.Activity{}
	activity.Expiration = now.Add(72 * time.Hour).Unix()
	activity.Type = model.UserActivation
	activity.Token = uuid.NewV4().String()
	activity.UserID = user.ID
	mail.SendMail([]string{user.Email}, activity.Token)
	return model.CreateActivity(activity)
}

func LoginUser(user *model.User) error {
	dbUser := model.FindUserByUserid(user.UserID)
	if !dbUser.Activated {
		return errors.New("Not activated")
	}
	if time.Now().After(time.Unix(dbUser.Expiration, 0)) {
		return errors.New("User expired")
	}
	err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	// TODO create token

	return errors.Wrap(err, "Login unsucessfull")
}

func ActivateUserByToken(tkn string) error {
	activity, err := model.FindActivityByToken(tkn)
	if err != nil {
		return errors.Wrap(err, "Token not valid")
	}
	if activity.Type != model.UserActivation {
		return errors.New("Token not valid")
	}
	now := time.Now().Unix()
	activity.Expiration = now
	activity.Used = now
	model.UpdateActivity(activity)
	u, err := model.FindUserByID(activity.UserID)
	if err != nil {
		return errors.Wrap(err, "Cannot find user")
	}
	u.Activated = true
	return model.SaveUser(u)
}

func ResetUserPassword(user *model.User) error {
	activity := &model.Activity{}
	activity.Expiration = time.Now().Add(72 * time.Hour).Unix()
	activity.Type = model.PasswordReset
	activity.Token = uuid.NewV4().String()
	activity.UserID = user.ID
	mail.SendMail([]string{user.Email}, activity.Token)
	return model.CreateActivity(activity)
}
