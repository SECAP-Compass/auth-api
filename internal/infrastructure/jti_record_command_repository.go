package infrastructure

import (
	"auth-api/internal/domain"
	"auth-api/internal/util"
	"context"

	"gorm.io/gorm"
)

type JtiRecordCommandRepository struct {
	db *gorm.DB
}

func NewJtiRecordCommandRepository(db *gorm.DB) domain.IJtiRecordCommandRepository {
	_ = db.AutoMigrate(&domain.JtiRecord{})
	return &JtiRecordCommandRepository{
		db: db,
	}
}

func (r *JtiRecordCommandRepository) Store(ctx context.Context, record *domain.JtiRecord) error {
	_, span := util.StartSpan(ctx)
	defer span.End()

	if err := r.db.Create(record).Error; err != nil {
		// r.db.Rollback()
		return err
	}

	return nil
}

func (r *JtiRecordCommandRepository) Delete(ctx context.Context, jti string) error {
	_, span := util.StartSpan(ctx)
	defer span.End()

	// Using Unscoped to delete the record permanently
	if err := r.db.Unscoped().Delete(&domain.JtiRecord{Id: jti}).Error; err != nil {
		// r.db.Rollback()
		return err
	}

	return nil
}
