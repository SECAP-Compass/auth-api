package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"auth-api/internal/application"
	"auth-api/internal/database"
	"auth-api/internal/infrastructure"
)

type Server struct {
	port int

	db           database.Service
	tokenService *application.TokenService
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	// TODO: Do I need this database service?
	dbService := database.New()
	if err := dbService.Health(); err != nil {
		panic(err)
	}

	db := dbService.GetDatabase()

	userRepository := infrastructure.NewUserRepository(db)
	jwtRepository := infrastructure.NewJtiRecordRepository(db)

	tokenService := application.NewTokenService(userRepository, jwtRepository)

	NewServer := &Server{
		port: port,

		db:           database.New(),
		tokenService: tokenService,
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
