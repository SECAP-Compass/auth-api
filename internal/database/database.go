package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Service interface {
	Health() error
	GetDatabase() *gorm.DB
}

type service struct {
	db *gorm.DB
}

var (
	database = os.Getenv("DB_DATABASE")
	password = os.Getenv("DB_PASSWORD")
	username = os.Getenv("DB_USERNAME")
	port     = os.Getenv("DB_PORT")
	host     = os.Getenv("DB_HOST")
)

func New() Service {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Istanbul", host, username, password, database, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	s := &service{db: db}
	return s
}

func (s *service) GetDatabase() *gorm.DB {
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

	d, err := s.db.DB()
	if err != nil {
		return err
	}

	if err := d.PingContext(ctx); err != nil {
		return err
	}

	return nil
}
