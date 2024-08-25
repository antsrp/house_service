package setup

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/antsrp/house_service/internal/repository/postgres"
	"github.com/antsrp/house_service/internal/rest"
	"github.com/antsrp/house_service/internal/service"
	"github.com/antsrp/house_service/pkg/config"
	"github.com/antsrp/house_service/pkg/crypt"
	ds "github.com/antsrp/house_service/pkg/infrastructure/db"
	rs "github.com/antsrp/house_service/pkg/infrastructure/rest"
	"github.com/antsrp/house_service/pkg/jwt"
	"github.com/antsrp/house_service/pkg/log"
)

func Setup(dbConfigPrefix string) (rest.Handler, rs.Settings, *postgres.Connection, *slog.Logger, error) {
	logger := log.NewSlogTextLogger()

	if err := config.Load(); err != nil {
		return rest.Handler{}, rs.Settings{}, nil, logger, fmt.Errorf("cannot load config data: %w", err)
	}

	key, err := os.ReadFile(".secret")
	if err != nil {
		return rest.Handler{}, rs.Settings{}, nil, logger, fmt.Errorf("cannot load jwt secret from file: %w", err)
	}

	srvSettings, err := config.Parse[rs.Settings]("SERVER")
	if err != nil {
		return rest.Handler{}, rs.Settings{}, nil, logger, fmt.Errorf("cannot parse server settings: %w", err)
	}

	dbSettings, err := config.Parse[ds.Settings](dbConfigPrefix)
	if err != nil {
		return rest.Handler{}, rs.Settings{}, nil, logger, fmt.Errorf("cannot parse database settings: %w", err)
	}

	dbConnection, err := postgres.NewConnection(context.Background(), dbSettings, logger)
	if err != nil {
		return rest.Handler{}, rs.Settings{}, nil, logger, fmt.Errorf("cannot init database connection: %w", err)
	}
	//defer dbConnection.Close()
	hs, fs, ss := postgres.NewHouseStorage(*dbConnection), postgres.NewFlatStorage(*dbConnection), postgres.NewSubscriberStorage(*dbConnection)

	subscriberService := service.NewSubscriberService(logger, ss)

	HFService := service.NewHouseFlatService(fs, hs, subscriberService)

	jwtService := jwt.NewJwtService(key)
	var cryptor crypt.Crypt
	tokenStorage := postgres.NewTokenStorage(*dbConnection)
	tokenService := service.NewTokenService(tokenStorage, jwtService)

	userStorage := postgres.NewUserStorage(*dbConnection)
	userService := service.NewUserService(tokenService, tokenService, userStorage, cryptor)

	h := rest.NewHandler(logger, srvSettings, HFService, userService, tokenService)

	return h, srvSettings, dbConnection, logger, nil
}
