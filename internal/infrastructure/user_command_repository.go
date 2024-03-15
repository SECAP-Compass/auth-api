package infrastructure

import (
	"auth-api/internal/domain"
	"context"

	"gorm.io/gorm"
)

type UserCommandRepository struct {
	db *gorm.DB
}

func NewUserCommandRepository(db *gorm.DB) domain.IUserCommandRepository {
	db.AutoMigrate(&domain.User{})
	return &UserCommandRepository{
		db: db,
	}
}

func (r *UserCommandRepository) Store(ctx context.Context, user *domain.User) error {
	if err := r.db.Create(user).Error; err != nil {
		// r.db.Rollback()
		return err
	}

	return nil
}

func (r *UserCommandRepository) Update(ctx context.Context, user *domain.User) error {
	if err := r.db.Save(user).Error; err != nil {
		// r.db.Rollback()
		return err
	}

	return nil
}

func (r *UserCommandRepository) Delete(ctx context.Context, id string) error {
	if err := r.db.Delete(&domain.User{}, id).Error; err != nil {
		// r.db.Rollback()
		return err
	}

	return nil
}
