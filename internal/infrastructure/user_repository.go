package infrastructure

import (
	"auth-api/internal/domain"

	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	db *mongo.Client
}

func NewUserRepository(db *mongo.Client) domain.IUserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	panic("implement")
}

func (r *UserRepository) FindByID(id string) (*domain.User, error) {
	panic("implement")
}

func (r *UserRepository) Store(user *domain.User) error {
	panic("implement")
}

func (r *UserRepository) Update(user *domain.User) error {
	panic("implement")
}

func (r *UserRepository) Delete(id string) error {
	panic("implement")
}
