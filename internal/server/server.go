package server

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"
	"github.com/lmittmann/tint"

	"auth-api/internal/application"
	"auth-api/internal/database"
	"auth-api/internal/infrastructure"
)

type Server struct {
	Port int
	App  *fiber.App

	db           database.Service
	tokenService *application.TokenService

	Validator *validator.Validate
}

func NewServer() *Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	app := fiber.New()
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

	return &Server{
		App:          app,
		Port:         port,
		db:           database.New(),
		tokenService: tokenService,
		Validator:    validator.New(validator.WithRequiredStructEnabled()),
	}
}
