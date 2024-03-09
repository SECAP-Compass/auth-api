package domain

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type Claims struct {
	jwt.StandardClaims
}

type JtiRecord struct {
	Id              string `json:"id" bson:"id"`
	Email           string `json:"email" bson:"email"`
	UserId          string `json:"userId" bson:"userId"`
	ExpireTimeStamp int64  `json:"expireTimeStamp" bson:"expireTimeStamp"`
}

// Do I need that?
type Jwt struct {
	Jwt *jwt.Token
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
			Audience:  "secap", // ?I think this should be Authority?
			Subject:   email,
		},
	), nil
}

func (Jwt *Jwt) String() string {
	tokenString, _ := Jwt.Jwt.SignedString([]byte("secret"))

	return tokenString
}

func (Jwt *Jwt) ToResponse() map[string]string {
	c := Jwt.Jwt.Claims.(*jwt.StandardClaims)

	return map[string]string{
		"access_token": Jwt.String(),
		"issuedAt":     strconv.FormatInt(c.IssuedAt, 10),
		"expiresAt":    strconv.FormatInt(c.ExpiresAt, 10),
	}
}

func NewJtiRecord(j *Jwt, userId string) *JtiRecord {
	c := j.Jwt.Claims.(*jwt.StandardClaims)

	return &JtiRecord{
		Id:              c.Id,
		Email:           c.Subject,
		UserId:          userId,
		ExpireTimeStamp: c.ExpiresAt,
	}
}
