package application

import (
	"auth-api/internal/domain"
	"context"
	"errors"
	"fmt"
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
	ctx = context.WithValue(ctx, "user", user)

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
		return nil, err
	}

	ctx = context.WithValue(ctx, "user", user)

	if !user.ComparePassword(r.Password) {
		return nil, fmt.Errorf("invalid password")
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
