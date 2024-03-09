package infrastructure

import (
	"auth-api/internal/domain"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	c *mongo.Collection
}

const USER_COLLECTION = "users"

func NewUserRepository(db *mongo.Database) domain.IUserRepository {
	c := db.Collection(USER_COLLECTION)
	return &UserRepository{
		c: c,
	}
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	filter := bson.M{"email": email}
	result := r.c.FindOne(ctx, filter)
	if result.Err() != nil {
		return nil, result.Err()
	}

	return r.parseUser(result)
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	filter := bson.M{"id": id}
	result := r.c.FindOne(ctx, filter)

	if result == nil {
		return nil, fmt.Errorf("user not found")
	}

	if result.Err() != nil {
		return nil, result.Err()
	}

	return r.parseUser(result)
}

func (r *UserRepository) Store(ctx context.Context, user *domain.User) error {
	_, err := r.c.InsertOne(ctx, user)
	return err
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	_, err := r.c.UpdateByID(ctx, user.ID, user)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	filter := bson.M{"id": id}
	_, err := r.c.DeleteOne(ctx, filter)
	return err
}

func (r *UserRepository) parseUser(result *mongo.SingleResult) (*domain.User, error) {
	user := &domain.User{}
	err := result.Decode(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
