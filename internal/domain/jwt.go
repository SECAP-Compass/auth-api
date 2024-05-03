package domain

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JtiRecord struct {
	gorm.Model `gorm:"primary_key:Id"`

	Id              string    `json:"id" gorm:"primary_key"`
	Email           string    `json:"email" gorm:"index"`
	UserId          uint      `json:"userId"`
	ExpireTimeStamp time.Time `json:"expireTimeStamp"`
}

type Jwt struct {
	Jwt *jwt.Token
}

type CustomClaims struct {
	*jwt.StandardClaims
	CityId    uint     `json:"cityId"`
	City      string   `json:"city"`
	Roles     []string `json:"roles"`
	Authority string   `json:"authority"`
}

const ISSUER = "secap-auth"

func NewJwt(email string) (*jwt.Token, error) {
	claims := &CustomClaims{
		&jwt.StandardClaims{
			Id:        uuid.NewString(),
			Issuer:    ISSUER,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Audience:  "secap", // ?I think this should be Authority?
			Subject:   email,
		},
		34, "istanbul",
		[]string{"buildingAdmin"},
		"istanbul",
	}

	return jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	), nil
}

func (j *Jwt) String() string {
	tokenString, _ := j.Jwt.SignedString([]byte("secret"))
	return tokenString
}

func (j *Jwt) ToResponse() map[string]string {
	c := j.Jwt.Claims.(*CustomClaims)
	return map[string]string{
		"access_token": j.String(),
		"issuedAt":     strconv.FormatInt(c.IssuedAt, 10),
		"expiresAt":    strconv.FormatInt(c.ExpiresAt, 10),
		"type":         "Bearer",
	}
}

func NewJtiRecord(j *Jwt, userId uint) *JtiRecord {
	c := j.Jwt.Claims.(*CustomClaims)

	return &JtiRecord{
		Id:              c.Id,
		Email:           c.Subject,
		UserId:          userId,
		ExpireTimeStamp: time.Unix(c.ExpiresAt, 0),
	}
}
