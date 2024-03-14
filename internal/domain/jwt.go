package domain

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Claims struct {
	jwt.StandardClaims
}

type JtiRecord struct {
	gorm.Model `gorm:"primary_key:Id"`

	Id              string    `json:"id" gorm:"primary_key"`
	Email           string    `json:"email"`
	UserId          uint      `json:"userId"`
	ExpireTimeStamp time.Time `json:"expireTimeStamp"`
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

func (j *Jwt) String() string {
	tokenString, _ := j.Jwt.SignedString([]byte("secret"))
	return tokenString
}

func (j *Jwt) ToResponse() map[string]string {
	c := j.Jwt.Claims.(*jwt.StandardClaims)

	return map[string]string{
		"access_token": j.String(),
		"issuedAt":     strconv.FormatInt(c.IssuedAt, 10),
		"expiresAt":    strconv.FormatInt(c.ExpiresAt, 10),
		"type":         "Bearer",
	}
}

func NewJtiRecord(j *Jwt, userId uint) *JtiRecord {
	c := j.Jwt.Claims.(*jwt.StandardClaims)

	return &JtiRecord{
		Id:              c.Id,
		Email:           c.Subject,
		UserId:          userId,
		ExpireTimeStamp: time.Unix(c.ExpiresAt, 0),
	}
}
