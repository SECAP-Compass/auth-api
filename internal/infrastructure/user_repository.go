package infrastructure

import (
	"auth-api/internal/domain"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

const USER_TABLE string = "users"

func NewUserRepository(db *gorm.DB) domain.IUserRepository {
	db.AutoMigrate(&domain.User{})
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := &domain.User{}

	if err := r.db.First(user, "email = ?", email).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return nil, fmt.Errorf("user.not.found.by.email")
		default:
			return nil, err
		}
	}

	return user, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id uint) (*domain.User, error) {
	user := &domain.User{}

	if err := r.db.First(user, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return nil, fmt.Errorf("user.not.found.by.id")
		default:
			return nil, err
		}
	}

	return user, nil
}

func (r *UserRepository) Store(ctx context.Context, user *domain.User) error {
	if err := r.db.Create(user).Error; err != nil {
		r.db.Rollback()
		return err
	}

	return nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	if err := r.db.Save(user).Error; err != nil {
		r.db.Rollback()
		return err
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	if err := r.db.Delete(&domain.User{}, id).Error; err != nil {
		r.db.Rollback()
		return err
	}

	return nil
}
