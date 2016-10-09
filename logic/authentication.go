package logic

import (
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
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

func LoginUser(user *model.User) (*model.User, error) {
	dbUser, err := model.FindUserByUserid(user.UserID)
	if !dbUser.Activated {
		return nil, errors.New("Not activated")
	}
	if time.Now().After(time.Unix(dbUser.Expiration, 0)) {
		return nil, errors.New("User expired")
	}
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		return nil, errors.Wrap(err, "Login unsucessfull")
	}
	return dbUser, nil
}

func ObtainToken(u *model.User) (string, error) {

	// Try to find existing
	// not expired token
	activeTkn, err := model.FindActiveByUser(u)
	if err == nil {
		return activeTkn.RefToken, err
	} else {
		log.Println("Cannot find active token %v \n", errors.Cause(err))
	}

	// No existing token found.
	// Create new one.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userid":  u.UserID,
		"email":   u.Email,
		"expires": time.Now().Add(72 * time.Hour).Unix(),
	})
	refTkn := uuid.NewV4().String()

	jwt, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}
	dbTkn := &model.Token{
		UserID:     u.ID,
		JwtToken:   jwt,
		RefToken:   refTkn,
		Expiration: time.Now().Add(72 * time.Hour).Unix(),
	}
	return refTkn, model.SaveToken(dbTkn)
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

func LogoutUser(tkn string) error {
	token, err := model.FindTokenByRef(tkn)
	if err != nil {
		return err
	}
	log.Printf("Obtained token %v\n", token)
	return model.InvalidateAllSessionByUsedID(token.UserID)
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

func generateNewToken(user *model.User, jwtToken string) model.Token {
	token := model.Token{}
	token.JwtToken = jwtToken
	token.UserID = user.ID
	token.RefToken = uuid.NewV4().String()
	token.Expiration = time.Now().Add(time.Hour * 72).Unix()
	return token
}
