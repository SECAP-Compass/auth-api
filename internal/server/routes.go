package server

import (
	"auth-api/internal/application"
	"auth-api/internal/util"
	"fmt"
	"log/slog"
	"net/http"

	validation "github.com/go-playground/validator/v10"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	json "github.com/json-iterator/go"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	// r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(TracerFn)

	r.Get("/health", s.healthHandler)
	r.Mount("/debug", middleware.Profiler())

	r.Post("/register", s.Register)
	r.Post("/login", s.Login)

	return r
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.db.Health())
	_, _ = w.Write(jsonResp)
}

// TODO: check status codes
// TODO: a conventional respone structure
func (s *Server) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	_, span := util.StartSpan(ctx)
	defer span.End()

	req := &application.UserRegisterRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = s.Validator.Struct(req); err != nil {
		errors := err.(validation.ValidationErrors)
		slog.Error("Error validating login request", slog.Any("error", errors))
		http.Error(w, fmt.Sprintf("validation errors: %s", errors), http.StatusBadRequest)
		return
	}

	jwt, err := s.tokenService.Register(ctx, req)
	if err != nil {
		slog.Error("Error registering user")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jwtByteArr, err := json.Marshal(jwt.ToResponse())
	if err != nil {
		slog.Error("Error marshalling jwt to response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Add("content-type", "application/json")
	_, _ = w.Write(jwtByteArr)
}

// TODO: check status codes
// TODO: a conventional respone structure
func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	_, span := util.StartSpan(ctx)
	defer span.End()

	req := &application.UserLoginRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()
	if err != nil {
		slog.Error("Error decoding login request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = s.Validator.Struct(req); err != nil {
		errors := err.(validation.ValidationErrors)
		slog.Error("Error validating login request", slog.Any("error", errors))
		http.Error(w, fmt.Sprintf("validation errors: %s", errors), http.StatusBadRequest)
		return
	}

	jwt, err := s.tokenService.Login(ctx, req)
	if err != nil {
		slog.Error("Error logging in user", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	jwtByteArr, err := json.Marshal(jwt.ToResponse())
	if err != nil {
		slog.Error("Error marshalling jwt to response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jwtByteArr)
}
