package service

import (
	"context"
	"time"

	"github.com/antsrp/house_service/internal/domain/models"
	"github.com/antsrp/house_service/internal/repository"
	"github.com/antsrp/house_service/pkg/jwt"
)

type TokenServicer interface {
	UserByToken(context.Context, string) (models.User, Error)
	TokenByID(context.Context, models.User) (string, Error)
	AddToken(context.Context, models.User, string) Error
}

type TokenCreator interface {
	CreateToken(context.Context, models.User) (string, Error)
}

type TokenService struct {
	tokenStorage repository.TokenStorage
	jwtService   jwt.Servicer
}

func NewTokenService(tokenStorage repository.TokenStorage, jwtService jwt.Servicer) TokenService {
	return TokenService{
		tokenStorage: tokenStorage,
		jwtService:   jwtService,
	}
}

func (t TokenService) UserByToken(ctx context.Context, token string) (models.User, Error) {
	user, err := t.tokenStorage.UserByToken(ctx, token)
	if err != nil {
		return models.User{}, NewServiceError(StatusByError(err), err.Cause(), DatabaseErrorCode)
	}

	return user, nil
}
func (t TokenService) TokenByID(ctx context.Context, user models.User) (string, Error) {
	token, err := t.tokenStorage.TokenByID(ctx, user)
	if err != nil {
		return "", NewServiceError(StatusByError(err), err.Cause(), DatabaseErrorCode)
	}

	return token, nil
}

func (t TokenService) AddToken(ctx context.Context, user models.User, token string) Error {
	if err := t.tokenStorage.AddToken(ctx, user, token, time.Now()); err != nil {
		return NewServiceError(StatusByError(err), err.Cause(), DatabaseErrorCode)
	}

	return nil
}

func (t TokenService) CreateToken(ctx context.Context, user models.User) (string, Error) {
	token, err := t.jwtService.NewToken(map[string]any{
		"id":         user.ID,
		"user_type":  user.UserType,
		"created_at": time.Now(),
	})
	if err != nil {
		return "", NewServiceError(Internal, err, CreateTokenErrorCode)
	}
	return token, nil
}

var _ TokenServicer = TokenService{}
