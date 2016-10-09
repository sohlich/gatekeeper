package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/pkg/errors"
	"github.com/sohlich/gatekeeper/logic"
	"github.com/sohlich/gatekeeper/model"
)

type UserHandler func(u *model.User, rw http.ResponseWriter)

func Register(rw http.ResponseWriter, req *http.Request) {
	UserBodyHandler(registerHandler)(rw, req)
}

func Login(rw http.ResponseWriter, req *http.Request) {
	UserBodyHandler(loginHandler)(rw, req)
}

func Logout(rw http.ResponseWriter, req *http.Request) {
	refTkn := req.Header.Get("token")
	err := logic.LogoutUser(refTkn)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

func ActivateUser(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	tkn := req.Form.Get("token")
	if len(tkn) == 0 {
		log.Println("cannot parse token")
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	if err := logic.ActivateUserByToken(tkn); err != nil {
		log.Println("cannot activate user: %s", err.Error())
		rw.WriteHeader(http.StatusNotFound)
		return
	}
}

func UserBodyHandler(h UserHandler) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		user := &model.User{}
		bytes, err := ioutil.ReadAll(req.Body)
		if err != nil {
			processError(rw, err)
			return
		}
		err = json.Unmarshal(bytes, user)
		if err != nil {
			processError(rw, err)
			return
		}
		h(user, rw)
	}
}

func loginHandler(u *model.User, rw http.ResponseWriter) {
	user, err := logic.LoginUser(u)
	if err != nil {
		log.Println("Login failed for user %s\n", u)
		rw.WriteHeader(http.StatusForbidden)
		return
	}
	tkn, err := logic.ObtainToken(user)
	if err != nil {
		log.Println("Cannot obtain token %+v \n", errors.Cause(err))
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.Header().Set("token", tkn)
	rw.WriteHeader(http.StatusOK)
}

func registerHandler(u *model.User, rw http.ResponseWriter) {
	logic.RegisterUser(u)
	json.NewEncoder(rw).Encode(u)
}

func processError(rw http.ResponseWriter, err error) {
	log.Printf("Cannot process request due: %+v\n", err)
	rw.WriteHeader(http.StatusInternalServerError)
}
