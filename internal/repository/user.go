package repository

import (
	"context"

	"github.com/antsrp/house_service/internal/domain/models"
)

type UserStorage interface {
	Add(context.Context, models.User) (models.User, DatabaseError)
	Get(context.Context, models.User) (models.User, DatabaseError)
}
