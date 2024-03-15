package infrastructure

import (
	"auth-api/internal/domain"
	"auth-api/internal/util"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type UserQueryRepository struct {
	db *gorm.DB
}

const USER_TABLE string = "users"

func NewUserQueryRepository(db *gorm.DB) domain.IUserQueryRepository {
	// 	db.AutoMigrate(&domain.User{})
	return &UserQueryRepository{
		db: db,
	}
}

func (r *UserQueryRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	_, span := util.StartSpan(ctx)
	defer span.End()

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

func (r *UserQueryRepository) FindByID(ctx context.Context, id uint) (*domain.User, error) {
	_, span := util.StartSpan(ctx)
	defer span.End()

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
