package repository

import (
	"context"

	"github.com/antsrp/house_service/internal/domain/models"
)

type HouseStorage interface {
	Create(context.Context, models.HouseCreateRequest) (models.HouseCreateResponse, DatabaseError)
	Flats(context.Context, models.HouseGetFlatsRequest, models.User) (models.HouseGetFlatsResponse, DatabaseError)
}
