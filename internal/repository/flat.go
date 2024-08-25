package repository

import (
	"context"

	"github.com/antsrp/house_service/internal/domain/models"
)

type FlatStorage interface {
	Create(context.Context, models.FlatCreateRequest) (models.FlatCreateResponse, DatabaseError)
	Update(context.Context, models.FlatUpdateRequest) (models.FlatUpdateResponse, DatabaseError)
	Get(context.Context, models.Flat) (models.Flat, DatabaseError)
}
