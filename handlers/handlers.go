package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/sohlich/attendence/logic"
	"github.com/sohlich/attendence/model"
)

type UserHandler func(u *model.User, rw http.ResponseWriter)

func Register(rw http.ResponseWriter, req *http.Request) {
	UserHandlerFunc(registerHandler)(rw, req)
}

func Login(rw http.ResponseWriter, req *http.Request) {
	UserHandlerFunc(loginHandler)(rw, req)
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

func UserHandlerFunc(h UserHandler) http.HandlerFunc {
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
	err := logic.LoginUser(u)
	if err != nil {
		log.Println("Login failed for user %s\n", u)
		rw.WriteHeader(http.StatusForbidden)
		return
	}
}

func registerHandler(u *model.User, rw http.ResponseWriter) {
	logic.RegisterUser(u)
	json.NewEncoder(rw).Encode(u)
}

func processError(rw http.ResponseWriter, err error) {
	log.Printf("Cannot process request due: %+v\n", err)
	rw.WriteHeader(http.StatusInternalServerError)
}
