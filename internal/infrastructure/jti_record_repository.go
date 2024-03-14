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

type JtiRecordRepository struct {
	db *gorm.DB
}

const JWT_COLLECTION string = "jtis"

func NewJtiRecordRepository(db *gorm.DB) domain.IJtiRecordRepository {
	db.AutoMigrate(&domain.JtiRecord{})
	return &JtiRecordRepository{
		db: db,
	}
}

func (r *JtiRecordRepository) FindByID(ctx context.Context, jti string) (*domain.JtiRecord, error) {
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

func (r *JtiRecordRepository) FindByUserID(ctx context.Context, userId uint) (*domain.JtiRecord, error) {
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

func (r *JtiRecordRepository) Store(ctx context.Context, record *domain.JtiRecord) error {
	if err := r.db.Create(record).Error; err != nil {
		r.db.Rollback()
		return err
	}

	return nil
}

func (r *JtiRecordRepository) Delete(ctx context.Context, jti string) error {

	// Using Unscoped to delete the record permanently
	if err := r.db.Unscoped().Delete(&domain.JtiRecord{Id: jti}).Error; err != nil {
		r.db.Rollback()
		return err
	}

	return nil
}
