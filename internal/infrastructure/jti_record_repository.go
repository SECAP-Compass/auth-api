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

const JWT_COLLECTION = "jtis"

func NewJtiRecordRepository(db *mongo.Database) domain.IJtiRecordRepository {
	c := db.Collection(JWT_COLLECTION)

	return &JtiRecordRepository{
		c: c,
	}
}

func (r *JtiRecordRepository) FindByID(ctx context.Context, jti string) (*domain.JtiRecord, error) {
	filter := bson.M{"_id": jti} // Is there a better way to do this? :)
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

func (r *JtiRecordRepository) FindByUserID(ctx context.Context, userId string) (*domain.JtiRecord, error) {
	filter := bson.M{"userId": userId}
	result := r.c.FindOne(ctx, filter)
	if result.Err() != nil {
		return nil, result.Err()
	}

	jti, err := r.parseJti(result)
	if err != nil {
		return nil, err
	}

	return jti, nil
}

func (r *JtiRecordRepository) Store(ctx context.Context, record *domain.JtiRecord) error {
	_, err := r.c.InsertOne(ctx, record)
	return err
}

func (r *JtiRecordRepository) Delete(ctx context.Context, jti string) error {
	filter := bson.M{"_id": jti}
	_, err := r.c.DeleteOne(ctx, filter)
	return err
}

func (r *JtiRecordRepository) parseJti(result *mongo.SingleResult) (*domain.JtiRecord, error) {
	jti := &domain.JtiRecord{}
	err := result.Decode(jti)
	if err != nil {
		return nil, err
	}

	return jti, nil
}
