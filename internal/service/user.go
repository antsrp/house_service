package service

import (
	"context"
	"errors"

	"github.com/antsrp/house_service/internal/domain/models"
	"github.com/antsrp/house_service/internal/repository"
	"github.com/antsrp/house_service/pkg/crypt"
)

type UserServicer interface {
	Login(context.Context, models.LoginRequest) (models.LoginResponse, Error)
	Register(context.Context, models.RegisterRequest) (models.RegisterResponse, Error)
}

type UserService struct {
	tokenService TokenServicer
	tokenCreator TokenCreator
	userStorage  repository.UserStorage
	cryptor      crypt.Cryptor
}

func NewUserService(tokenService TokenServicer, tokenCreator TokenCreator, userStorage repository.UserStorage, cryptor crypt.Cryptor) UserService {
	return UserService{
		tokenService: tokenService,
		tokenCreator: tokenCreator,
		userStorage:  userStorage,
		cryptor:      cryptor,
	}
}

func (u UserService) Login(ctx context.Context, req models.LoginRequest) (models.LoginResponse, Error) {
	hash := u.cryptor.Hash(*req.Password)
	user := models.User{
		ID:       *req.UserID,
		Password: hash,
	}
	var err repository.DatabaseError
	if user, err = u.userStorage.Get(ctx, user); err != nil {
		if errors.Is(err.Cause(), repository.ErrEntityNotFound) {
			return models.LoginResponse{}, NewServiceError(StatusByError(err), ErrUserNotFound, DatabaseErrorCode)
		}
		return models.LoginResponse{}, NewServiceError(StatusByError(err), err.Cause(), DatabaseErrorCode)
	}
	token, terr := u.tokenService.TokenByID(ctx, user)
	if terr != nil {
		if !errors.Is(terr.Cause(), repository.ErrEntityNotFound) {
			return models.LoginResponse{}, terr
		}
		if token, terr = u.tokenCreator.CreateToken(ctx, user); terr != nil {
			return models.LoginResponse{}, terr
		}
		if err := u.tokenService.AddToken(ctx, user, token); err != nil {
			return models.LoginResponse{}, err
		}
	}
	return models.LoginResponse{
		DummyLoginResponse: models.DummyLoginResponse{
			Token: token,
		},
	}, nil
}
func (u UserService) Register(ctx context.Context, req models.RegisterRequest) (models.RegisterResponse, Error) {
	hash := u.cryptor.Hash(*req.Password)
	user := models.User{
		Email:    *req.Email,
		Password: hash,
		UserType: *req.UserType,
	}
	var err repository.DatabaseError
	if user, err = u.userStorage.Add(ctx, user); err != nil {
		if errors.Is(err.Cause(), repository.ErrEntityAlreadyExists) {
			return models.RegisterResponse{}, NewServiceError(StatusByError(err), ErrUserAlreadyExists, DatabaseErrorCode)
		}
		return models.RegisterResponse{}, NewServiceError(StatusByError(err), err.Cause(), DatabaseErrorCode)
	}

	return models.RegisterResponse{UserID: user.ID}, nil
}

var _ UserServicer = UserService{}
