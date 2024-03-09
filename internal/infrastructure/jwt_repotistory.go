package infrastructure

import (
	"auth-api/internal/domain"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/mongo"
)

// TODO:  Delete expired jwts?
type JwtRepository struct {
	db *mongo.Client
}

func NewJwtRepository(db *mongo.Client) domain.IJwtRepository {
	return &JwtRepository{
		db: db,
	}
}

func (r *JwtRepository) FindByID(jti string) (*jwt.Token, error) {
	panic("not implemented")
}

func (r *JwtRepository) Store(jwt *jwt.Token) error {
	panic("not implemented")
}
