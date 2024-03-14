package domain

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Email     string `json:"email" gorm:"index:unique"`
	Password  string `json:"password" `
	Authority string `json:"authority"`
}

func NewUser(email, password, authority string) *User {

	cryptedPassword := bcrpytPassword(password)
	return &User{
		Email:     email,
		Password:  cryptedPassword,
		Authority: authority,
	}
}

func (u *User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func bcrpytPassword(password string) string {
	cryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		panic(err)
	}
	return string(cryptedPassword)
}
