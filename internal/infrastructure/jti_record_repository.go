package infrastructure

import (
	"auth-api/internal/domain"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// TODO:  Delete expired jwts?
type JtiRecordRepository struct {
	c *mongo.Collection
}

const JWT_COLLECTION = "jwts"

func NewJtiRecordRepository(db *mongo.Database) domain.IJtiRecordRepository {
	c := db.Collection(JWT_COLLECTION)

	return &JtiRecordRepository{
		c: c,
	}
}

func (r *JtiRecordRepository) FindByID(ctx context.Context, jti string) (*domain.JtiRecord, error) {
	filter := bson.M{"id": jti} // Is there a better way to do this? :)
	result := r.c.FindOne(ctx, filter)
	if result.Err() != nil {
		return nil, result.Err()
	}

	jwt := &domain.JtiRecord{}
	err := result.Decode(jwt)
	if err != nil {
		return nil, err
	}

	return jwt, nil
}

func (r *JtiRecordRepository) Store(ctx context.Context, record *domain.JtiRecord) error {
	_, err := r.c.InsertOne(ctx, record)
	return err
}
