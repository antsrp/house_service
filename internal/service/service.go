package service

import (
	"context"
	"fmt"

	"github.com/antsrp/house_service/internal/domain/models"
	"github.com/antsrp/house_service/internal/repository"
)

type HouseFlatServicer interface {
	CreateHouse(context.Context, models.HouseCreateRequest) (models.HouseCreateResponse, Error)
	CreateFlat(context.Context, models.FlatCreateRequest) (models.FlatCreateResponse, Error)
	UpdateFlat(context.Context, models.FlatUpdateRequest) (models.FlatUpdateResponse, Error)
	Flats(context.Context, models.HouseGetFlatsRequest, models.User) (models.HouseGetFlatsResponse, Error)
	AddSubscriber(context.Context, string, int) Error
}

type HouseFlatService struct {
	flatStorage       repository.FlatStorage
	houseStorage      repository.HouseStorage
	subscriberService SubscriberServicer
}

var _ HouseFlatServicer = HouseFlatService{}

func NewHouseFlatService(fs repository.FlatStorage, hs repository.HouseStorage, ss SubscriberServicer) HouseFlatService {
	return HouseFlatService{
		flatStorage:       fs,
		houseStorage:      hs,
		subscriberService: ss,
	}
}

func (h HouseFlatService) CreateHouse(ctx context.Context, req models.HouseCreateRequest) (models.HouseCreateResponse, Error) {
	house, err := h.houseStorage.Create(ctx, req)
	if err != nil {
		return models.HouseCreateResponse{}, NewServiceError(StatusByError(err), err.Cause(), DatabaseErrorCode)
	}

	return house, nil
}
func (h HouseFlatService) CreateFlat(ctx context.Context, req models.FlatCreateRequest) (models.FlatCreateResponse, Error) {
	flat, err := h.flatStorage.Create(ctx, req)
	if err != nil {
		return models.FlatCreateResponse{}, NewServiceError(StatusByError(err), err.Cause(), DatabaseErrorCode)
	}

	go h.subscriberService.SendEmail(ctx, flat.HouseID)

	return flat, nil
}
func (h HouseFlatService) UpdateFlat(ctx context.Context, req models.FlatUpdateRequest) (models.FlatUpdateResponse, Error) {
	flat, err := h.flatStorage.Get(ctx, models.Flat{ID: req.ID})
	if err != nil {
		return models.FlatUpdateResponse{}, NewServiceError(StatusByError(err), fmt.Errorf("cannot get status of flat: %w", err.Cause()), DatabaseErrorCode)
	}
	if flat.Status == models.OnModeration { // already moderated
		return models.FlatUpdateResponse{}, NewServiceError(Conflict, fmt.Errorf("cannot update flat: already moderated"), DatabaseErrorCode)
	}

	stat := models.OnModeration
	if _, err := h.flatStorage.Update(ctx, models.FlatUpdateRequest{ID: req.ID, Status: &stat}); err != nil {
		return models.FlatUpdateResponse{}, NewServiceError(StatusByError(err), fmt.Errorf("cannot update flat status to `on moderate`: %w", err.Cause()), DatabaseErrorCode)
	}

	if req.Status == nil { // if status is not set, change it to approved
		stat = models.Approved
		req.Status = &stat
	}

	res, err := h.flatStorage.Update(ctx, req)
	if err != nil {
		return models.FlatUpdateResponse{}, NewServiceError(StatusByError(err), err.Cause(), DatabaseErrorCode)
	}

	return res, nil
}
func (h HouseFlatService) Flats(ctx context.Context, req models.HouseGetFlatsRequest, user models.User) (models.HouseGetFlatsResponse, Error) {
	flats, err := h.houseStorage.Flats(ctx, req, user)
	if err != nil {
		return models.HouseGetFlatsResponse{}, NewServiceError(StatusByError(err), err.Cause(), DatabaseErrorCode)
	}

	return flats, nil
}

func (h HouseFlatService) AddSubscriber(ctx context.Context, email string, id int) Error {
	if err := h.subscriberService.Add(ctx, email, id); err != nil {
		return err
	}
	return nil
}
