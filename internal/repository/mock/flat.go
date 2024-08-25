package mock

import (
	"context"

	"github.com/antsrp/house_service/internal/domain/models"
	"github.com/antsrp/house_service/internal/repository"
)

type FlatStorage struct {
	base *Base
}

var _ repository.FlatStorage = FlatStorage{}

func NewFlatStorage(base *Base) FlatStorage {
	return FlatStorage{
		base: base,
	}
}

func (f FlatStorage) Create(ctx context.Context, req models.FlatCreateRequest) (models.FlatCreateResponse, repository.DatabaseError) {
	if _, err := f.base.GetHouse(req.HouseID); err != nil {
		return models.FlatCreateResponse{}, NewMockError(false, err)
	}
	flat := models.Flat{
		HouseID: req.HouseID,
		Price:   *req.Price,
		Room:    req.Room,
		Status:  models.Created,
	}
	id := f.base.AddFlat(flat)
	flat.ID = id

	return models.FlatCreateResponse{
		Flat: flat,
	}, nil
}

func (f FlatStorage) Get(ctx context.Context, req models.Flat) (models.Flat, repository.DatabaseError) {
	flat, err := f.base.GetFlat(req.ID)
	if err != nil {
		return models.Flat{}, NewMockError(false, err)
	}
	return flat, nil
}

func (f FlatStorage) Update(ctx context.Context, req models.FlatUpdateRequest) (models.FlatUpdateResponse, repository.DatabaseError) {
	flat := models.Flat{
		ID:    req.ID,
		Price: *req.Price,
		Room:  req.Room,
	}
	if req.Status != nil {
		flat.Status = *req.Status
	}
	if err := f.base.UpdateFlat(flat); err != nil {
		return models.FlatUpdateResponse{}, NewMockError(false, err)
	}

	return models.FlatUpdateResponse{
		Flat: flat,
	}, nil
}
