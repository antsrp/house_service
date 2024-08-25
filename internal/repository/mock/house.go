package mock

import (
	"context"
	"time"

	"github.com/antsrp/house_service/internal/domain/models"
	"github.com/antsrp/house_service/internal/repository"
)

type HouseStorage struct {
	base *Base
}

var _ repository.HouseStorage = HouseStorage{}

func NewHouseStorage(base *Base) HouseStorage {
	return HouseStorage{
		base: base,
	}
}

func (f HouseStorage) Create(ctx context.Context, req models.HouseCreateRequest) (models.HouseCreateResponse, repository.DatabaseError) {
	house := models.House{
		Address:   req.Address,
		Year:      *req.Year,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if req.Developer != nil {
		house.Developer = *req.Developer
	}

	id := f.base.AddHouse(house)
	house.ID = id

	return models.HouseCreateResponse{
		House: house,
	}, nil
}

func (f HouseStorage) Flats(ctx context.Context, req models.HouseGetFlatsRequest, user models.User) (models.HouseGetFlatsResponse, repository.DatabaseError) {
	return f.base.Flats(req, user)
}
