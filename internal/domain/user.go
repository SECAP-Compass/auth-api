package domain

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string `json:"id" protobuf:"bytes,1,opt,name=id"`
	Email     string `json:"email" protobuf:"bytes,2,opt,name=email"`
	Password  string `json:"password" protobuf:"bytes,3,opt,name=password"`
	Authority string `json:"authority" protobuf:"bytes,5,opt,name=authority"`
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
	cryptedPassword := bcrpytPassword(password)
	return u.Password == cryptedPassword
}

func bcrpytPassword(password string) string {
	cryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		panic(err)
	}
	return string(cryptedPassword)
}
