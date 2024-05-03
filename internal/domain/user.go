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
	CityId    uint   `json:"cityId"`
}

func NewUser(email, password, authority string, cityId uint) *User {
	encryptedPassword := bcrpytPassword(password)
	return &User{
		Email:     email,
		Password:  encryptedPassword,
		Authority: authority,
		CityId:    cityId,
	}
}

func (u *User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func bcrpytPassword(password string) string {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		panic(err)
	}
	return string(encryptedPassword)
}
