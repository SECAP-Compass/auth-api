package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service interface {
	Health() error
	GetClient() *mongo.Client
	GetDatabase() *mongo.Database
}

type service struct {
	db     *mongo.Database
	client *mongo.Client
}

var (
	host     = os.Getenv("DB_HOST")
	port     = os.Getenv("DB_PORT")
	database = os.Getenv("DB_DATABASE")
)

func New() Service {
	client, err := mongo.Connect(
		context.Background(),
		options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", host, port)),
		options.Client().SetAuth(options.Credential{
			Username: os.Getenv("DB_USERNAME"),
			Password: os.Getenv("DB_PASSWORD"),
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	return &service{
		db:     client.Database(database),
		client: client,
	}
}

func (s *service) GetClient() *mongo.Client {
	return s.client
}

func (s *service) GetDatabase() *mongo.Database {
	return s.db
}

func (s *service) Health() error {

	healthy := false
	for i := 0; !healthy && i < 5; i++ {
		err := s.ping()
		if err == nil {
			healthy = true
		}
	}
	if !healthy {
		return fmt.Errorf("could not connect to database")
	}

	return nil
}

func (s *service) ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}
