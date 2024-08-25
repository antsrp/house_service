package repository

import (
	"context"
	"time"

	"github.com/antsrp/house_service/internal/domain/models"
)

type TokenStorage interface {
	UserByToken(context.Context, string) (models.User, DatabaseError)
	TokenByID(context.Context, models.User) (string, DatabaseError)
	AddToken(context.Context, models.User, string, time.Time) DatabaseError
}
