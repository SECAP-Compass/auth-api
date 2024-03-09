package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"auth-api/internal/database"
	"auth-api/internal/domain"
	"auth-api/internal/infrastructure"
)

type Server struct {
	port int

	db             database.Service
	userRepository domain.IUserRepository
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	dbService := database.New()

	userRepository := infrastructure.NewUserRepository(dbService.GetClient())

	NewServer := &Server{
		port: port,

		db:             database.New(),
		userRepository: userRepository,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
