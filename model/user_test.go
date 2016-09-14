package model

import "testing"
import "golang.org/x/crypto/bcrypt"

func TestBcrypt(t *testing.T) {
	password := "password"
	pass, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	bPass := string(pass)
	err := bcrypt.CompareHashAndPassword([]byte(bPass), []byte(password))
	if err != nil {
		t.FailNow()
	}
}
