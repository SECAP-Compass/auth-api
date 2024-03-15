package infrastructure

import (
	"auth-api/internal/domain"
	"context"
	"fmt"

	"gorm.io/gorm"
)

/*
 - With this repository implementation, we generate SSO someway.
*/

type JtiRecordQueryRepository struct {
	db *gorm.DB
}

func NewJtiRecordQueryRepository(db *gorm.DB) domain.IJtiRecordQueryRepository {
	// db.AutoMigrate(&domain.JtiRecord{})
	return &JtiRecordQueryRepository{
		db: db,
	}
}

func (r *JtiRecordQueryRepository) FindByID(ctx context.Context, jti string) (*domain.JtiRecord, error) {
	jtiRecord := &domain.JtiRecord{}

	if err := r.db.First(jtiRecord, jti).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return nil, fmt.Errorf("jti.not.found.by.id")
		default:
			return nil, err
		}
	}

	return jtiRecord, nil
}

func (r *JtiRecordQueryRepository) FindByUserID(ctx context.Context, userId uint) (*domain.JtiRecord, error) {
	jtiRecord := &domain.JtiRecord{}

	if err := r.db.First(jtiRecord, "user_id = ?", userId).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return nil, fmt.Errorf("jti.not.found.by.userId")
		default:
			return nil, err
		}
	}

	return jtiRecord, nil
}
