package mock

import (
	"github.com/antsrp/house_service/internal/domain/models"
	"github.com/antsrp/house_service/internal/repository"
)

type Base struct {
	flats     map[int]models.Flat
	houses    map[int]models.House
	cntFlats  int
	cntHouses int
}

func NewBase() Base {
	return Base{
		flats:  make(map[int]models.Flat),
		houses: make(map[int]models.House),
	}
}

func (b *Base) AddHouse(house models.House) int {
	b.houses[b.cntHouses+1] = house
	b.cntHouses++
	return b.cntHouses
}

func (b *Base) AddFlat(flat models.Flat) int {
	b.flats[b.cntFlats+1] = flat
	b.cntFlats++
	return b.cntFlats
}

func (b Base) Flats(req models.HouseGetFlatsRequest, user models.User) (models.HouseGetFlatsResponse, repository.DatabaseError) {
	house, err := b.GetHouse(req.ID)
	if err != nil {
		return models.HouseGetFlatsResponse{}, NewMockError(false, err)
	}

	var flats []models.Flat

	for k, v := range b.flats {
		if v.HouseID != req.ID {
			continue
		}
		if user.UserType == models.Moderator || (v.Status == models.Approved) {
			v.ID = k
			flats = append(flats, v)
		}
	}

	return models.HouseGetFlatsResponse{
		House: house,
		Flats: flats,
	}, nil
}

func (b Base) GetHouse(id int) (models.House, error) {
	house, found := b.houses[id]
	if !found {
		return models.House{}, repository.ErrEntityNotFound
	}
	house.ID = id

	return house, nil
}

func (b Base) GetFlat(id int) (models.Flat, error) {
	flat, found := b.flats[id]
	if !found {
		return models.Flat{}, repository.ErrEntityNotFound
	}
	flat.ID = id

	return flat, nil
}

func (b Base) UpdateFlat(flat models.Flat) error {
	if _, found := b.flats[flat.ID]; !found {
		return repository.ErrEntityNotFound
	}
	b.flats[flat.ID] = flat

	return nil
}
