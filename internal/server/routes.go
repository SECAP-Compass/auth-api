package server

import (
	"auth-api/internal/application"
	"log"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	json "github.com/json-iterator/go"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", s.HelloWorldHandler)
	r.Get("/health", s.healthHandler)

	r.Post("/register", s.register)
	r.Post("/login", s.login)

	return r
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.db.Health())
	_, _ = w.Write(jsonResp)
}

func (s *Server) register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &application.UserRegisterRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := &application.UserLoginRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()
	if err != nil {
		slog.Error("Error decoding login request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jwt, err := s.tokenService.Login(ctx, req)
	if err != nil {
		slog.Error("Error logging in user")
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
