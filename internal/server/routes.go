package server

import (
	"auth-api/internal/domain"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-jwt/jwt"
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
	var req *UserRegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := domain.NewUser(req.Email, req.Password, req.Authority)

	jwt, err := s.generateJwt(user.Email)
	if err != nil {
		http.Error(w, "could not generate jwt"+err.Error(), http.StatusInternalServerError)
		return
	}

	if err = s.userRepository.Store(user); err != nil {
		http.Error(w, "could not store user"+err.Error(), http.StatusInternalServerError)
		return

	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Add("id", user.ID)

	jwtByteArr, err := json.Marshal(jwt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(jwtByteArr)
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	var req *UserLoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := s.userRepository.FindByEmail(req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if !user.ComparePassword(req.Password) {
		http.Error(w, "invalid password", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)

	jwt, err := s.generateJwt(user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jwtByteArr, err := json.Marshal(jwt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(jwtByteArr)
}

func (s *Server) generateJwt(email string) (*jwt.Token, error) {
	jwt, err := domain.NewJwt(email)
	if err != nil {
		return nil, err
	}

	err = s.jwtRepository.Store(jwt)
	if err != nil {
		return nil, err
	}

	return jwt, nil
}
