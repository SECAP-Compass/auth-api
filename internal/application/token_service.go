package application

import (
	"auth-api/internal/domain"
	"auth-api/internal/util"
	"context"
	"errors"
	"fmt"
	"log/slog"
)

type TokenService struct {
	userQueryRepository   domain.IUserQueryRepository
	userCommandRepository domain.IUserCommandRepository

	JtiRecordQueryRepository   domain.IJtiRecordQueryRepository
	JtiRecordCommandRepository domain.IJtiRecordCommandRepository
}

func NewTokenService(
	userQueryRepository domain.IUserQueryRepository,
	userCommandRepository domain.IUserCommandRepository,
	jtiRecordQueryRepository domain.IJtiRecordQueryRepository,
	jtiRecordCommandRepository domain.IJtiRecordCommandRepository,
) *TokenService {
	return &TokenService{
		userQueryRepository:        userQueryRepository,
		userCommandRepository:      userCommandRepository,
		JtiRecordQueryRepository:   jtiRecordQueryRepository,
		JtiRecordCommandRepository: jtiRecordCommandRepository,
	}
}

func (s *TokenService) Register(ctx context.Context, r *UserRegisterRequest) (map[string]string, error) {
	user := domain.NewUser(r.Email, r.Password, r.Authority, r.CityId)

	if _, err := s.userQueryRepository.FindByEmail(ctx, r.Email); err == nil {
		return nil, fmt.Errorf("user already exists with email %s", r.Email)
	}

	if err := s.userCommandRepository.Store(ctx, user); err != nil {
		return nil, err
	}
	// nolint: staticcheck
	ctx = context.WithValue(ctx, "user", user)

	jwt, err := s.generateJwt(ctx)
	if err != nil {
		return nil, err
	}

	if err = s.saveJtiRecord(ctx, jwt); err != nil {
		return nil, err
	}

	return jwt.ToResponse(), nil
}

func (s *TokenService) Login(ctx context.Context, r *UserLoginRequest) (map[string]string, error) {
	ctx, span := util.StartSpan(ctx)
	defer span.End()

	user, err := s.userQueryRepository.FindByEmail(ctx, r.Email)
	if err != nil {
		slog.Error("Error finding user by email", slog.String("email", r.Email))
		return nil, err
	}

	// nolint: staticcheck
	ctx = context.WithValue(ctx, "user", user)

	if !user.ComparePassword(r.Password) {
		return nil, fmt.Errorf("invalid password")
	}

	jti, err := s.JtiRecordQueryRepository.FindByUserID(ctx, user.ID)
	if err == nil { // If there is no error, there is a jti record for this user
		go func() {
			if err := s.JtiRecordCommandRepository.Delete(ctx, jti.Id); err != nil {
				slog.Error("Error deleting jti record", slog.Any("jtiRecord", jti))
			}
		}()
	}

	jwt, err := s.generateJwt(ctx)
	if err != nil {
		return nil, err
	}

	if err = s.saveJtiRecord(ctx, jwt); err != nil {
		return nil, err
	}
	return jwt.ToResponse(), nil
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

	return &domain.Jwt{
		Jwt: jwt,
	}, nil
}

func (s *TokenService) saveJtiRecord(ctx context.Context, jwt *domain.Jwt) error {
	u := ctx.Value("user").(*domain.User)
	if u == nil {
		return errors.New("user not found in context")
	}

	jtiRecord := domain.NewJtiRecord(jwt, u.ID)
	err := s.JtiRecordCommandRepository.Store(ctx, jtiRecord)
	if err != nil {
		slog.Error("Error storing jti record", slog.Any("jtiRecord", jtiRecord))
		return err
	}

	return nil
}
