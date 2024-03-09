package domain

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type Claims struct {
	jwt.StandardClaims
}

const ISSUER = "secap-auth"

func NewJwt(email string) (*jwt.Token, error) {
	return jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		&jwt.StandardClaims{
			Id:        uuid.NewString(),
			Issuer:    ISSUER,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Audience:  "secap",
			Subject:   email,
		},
	), nil
}
