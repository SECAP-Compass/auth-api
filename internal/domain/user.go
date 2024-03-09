package domain

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string `json:"id" bson:"id"`
	Email     string `json:"email" bson:"email"`
	Password  string `json:"password" bson:"password"`
	Authority string `json:"authority" bson:"authority"`
}

func NewUser(email, password, authority string) *User {
	id := uuid.New().String()

	cryptedPassword := bcrpytPassword(password)
	return &User{
		ID:        id,
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
