package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"
	"github.com/lmittmann/tint"

	"auth-api/internal/application"
	"auth-api/internal/database"
	"auth-api/internal/infrastructure"
)

type Server struct {
	port int

	db           database.Service
	tokenService *application.TokenService

	Validator *validator.Validate
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	initTracer()

	// TODO: Do I need this database service?
	dbService := database.New()
	if err := dbService.Health(); err != nil {
		panic(err)
	}

	// Persistence layer
	db := dbService.GetDatabase()

	userQueryRepository := infrastructure.NewUserQueryRepository(db)
	userCommandRepository := infrastructure.NewUserCommandRepository(db)
	jtiRecordQueryRepository := infrastructure.NewJtiRecordQueryRepository(db)
	jtiRecordCommandRepository := infrastructure.NewJtiRecordCommandRepository(db)

	// Application layer
	tokenService := application.NewTokenService(userQueryRepository, userCommandRepository, jtiRecordQueryRepository, jtiRecordCommandRepository)

	// Logger
	loggerHandler := tint.NewHandler(os.Stdout, &tint.Options{AddSource: true})
	logger := slog.New(loggerHandler)
	slog.SetDefault(logger)

	NewServer := &Server{
		port:         port,
		db:           database.New(),
		tokenService: tokenService,
		Validator:    validator.New(validator.WithRequiredStructEnabled()),
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
