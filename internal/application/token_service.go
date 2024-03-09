package application

import (
	"auth-api/internal/domain"
	"context"
	"errors"
	"fmt"
	"log/slog"
)

type TokenService struct {
	userRepository domain.IUserRepository
	jtiRepository  domain.IJtiRecordRepository
}

func NewTokenService(userRepository domain.IUserRepository, jtiRepository domain.IJtiRecordRepository) *TokenService {
	return &TokenService{
		userRepository: userRepository,
		jtiRepository:  jtiRepository,
	}
}

func (s *TokenService) Register(ctx context.Context, r *UserRegisterRequest) (*domain.Jwt, error) {
	user := domain.NewUser(r.Email, r.Password, r.Authority)

	// nolint: staticcheck
	ctx = context.WithValue(ctx, "user", user)

	if _, err := s.userRepository.FindByEmail(ctx, r.Email); err == nil {
		return nil, fmt.Errorf("user already exists with email %s", r.Email)
	}

	jwt, err := s.generateJwt(ctx)
	if err != nil {
		return nil, err
	}

	if err = s.saveJtiRecord(ctx, jwt); err != nil {
		return nil, err
	}

	if err = s.userRepository.Store(ctx, user); err != nil {
		return nil, err
	}

	return jwt, nil
}

func (s *TokenService) Login(ctx context.Context, r *UserLoginRequest) (*domain.Jwt, error) {
	user, err := s.userRepository.FindByEmail(ctx, r.Email)
	if err != nil {
		slog.Error("Error finding user by email", slog.String("email", r.Email))
		return nil, err
	}

	// nolint: staticcheck
	ctx = context.WithValue(ctx, "user", user)

	if !user.ComparePassword(r.Password) {
		return nil, fmt.Errorf("invalid password")
	}

	jti, err := s.jtiRepository.FindByUserID(ctx, user.ID)
	if err == nil { // If there is no error, there is a jti record for this user
		if err := s.jtiRepository.Delete(ctx, jti.Id); err != nil {
			return nil, err
		}
	}

	jwt, err := s.generateJwt(ctx)
	if err != nil {
		return nil, err
	}

	if err = s.saveJtiRecord(ctx, jwt); err != nil {
		return nil, err
	}

	return jwt, nil
}

func (s *TokenService) generateJwt(ctx context.Context) (*domain.Jwt, error) {
	u := ctx.Value("user").(*domain.User)
	if u == nil {
		return nil, errors.New("user not found in context")
	}

	jwt, err := domain.NewJwt(u.Email)
	if err != nil {
		return nil, err
	}

	return &domain.Jwt{Jwt: jwt}, nil
}

func (s *TokenService) saveJtiRecord(ctx context.Context, jwt *domain.Jwt) error {
	u := ctx.Value("user").(*domain.User)
	if u == nil {
		return errors.New("user not found in context")
	}
	jtiRecord := domain.NewJtiRecord(jwt, u.ID)
	err := s.jtiRepository.Store(ctx, jtiRecord)
	if err != nil {
		return err
	}

	return nil
}
