package server

import (
	"auth-api/internal/application"
	"auth-api/internal/util"
	"fmt"
	"log/slog"
	"net/http"

	validation "github.com/go-playground/validator/v10"
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
)

func (s *Server) RegisterRoutes() {
	r := s.App
	r.Use(otelfiber.Middleware(otelfiber.WithPropagators(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))))
	r.Use(logger.New())

	r.Post("/register", s.Register)
	r.Post("/login", s.Login)

}

// TODO: check status codes
// TODO: a conventional respone structure
func (s *Server) Register(c *fiber.Ctx) error {
	ctx := c.UserContext()

	ctx, span := util.StartSpan(ctx)
	defer span.End()

	req := &application.UserRegisterRequest{}

	if err := c.BodyParser(req); err != nil {
		slog.Error("Error decoding register request")
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	if err := s.Validator.Struct(req); err != nil {
		errors := err.(validation.ValidationErrors)
		slog.Error("Error validating login request", slog.Any("error", errors))
		return c.Status(http.StatusBadRequest).SendString(fmt.Sprintf("validation errors: %s", errors))
	}

	jwt, err := s.tokenService.Register(ctx, req)
	if err != nil {
		slog.Error("Error registering user")
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(201).JSON(jwt.ToResponse())
}

// TODO: check status codes
// TODO: a conventional respone structure
func (s *Server) Login(c *fiber.Ctx) error {
	ctx := c.UserContext()

	otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(c.GetReqHeaders()))
	ctx, span := util.StartSpan(ctx)
	defer span.End()

	span.SetAttributes(attribute.String("login", "login"))
	req := &application.UserLoginRequest{}

	if err := c.BodyParser(req); err != nil {
		slog.Error("Error decoding login request")
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	if err := s.Validator.Struct(req); err != nil {
		errors := err.(validation.ValidationErrors)
		slog.Error("Error validating login request", slog.Any("error", errors))
		return c.Status(http.StatusBadRequest).SendString(fmt.Sprintf("validation errors: %s", errors))
	}

	jwt, err := s.tokenService.Login(ctx, req)
	if err != nil {
		slog.Error("Error logging in user", slog.Any("error", err))
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(200).JSON(jwt.ToResponse())
}
