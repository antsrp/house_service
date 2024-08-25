package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/antsrp/house_service/internal/domain/models"
	"github.com/antsrp/house_service/internal/repository"
	"github.com/antsrp/house_service/internal/repository/mock"
	"github.com/antsrp/house_service/internal/service"
	"github.com/stretchr/testify/require"
)

var (
	srv     service.HouseFlatServicer
	houseID int
	flats   []models.Flat
)

func TestMain(m *testing.M) {
	base := mock.NewBase()
	fs, hs, ss := mock.NewFlatStorage(&base), mock.NewHouseStorage(&base), service.NewMockSubscriberService()
	srv = service.NewHouseFlatService(fs, hs, ss)

	for i := 0; i < 2; i++ {
		flats = append(flats, models.Flat{ID: i + 1, HouseID: 1, Price: 123 + i, Room: 4 - i, Status: models.Created})
	}

	m.Run()
}

func TestAddHouses(t *testing.T) {
	year := 2023
	developer := `PEEK`
	resp, err := srv.CreateHouse(context.Background(), models.HouseCreateRequest{Address: "addr1", Year: &year, Developer: &developer})
	if err != nil {
		require.Equalf(t, true, false, "expected no error, actual %s", err.Cause())
	}

	houseID = resp.ID
}

func TestAddFlatNonExistingHouse(t *testing.T) {
	price := 1213
	_, err := srv.CreateFlat(context.Background(), models.FlatCreateRequest{HouseID: 2, Price: &price, Room: 2})
	if !errors.Is(err.Cause(), repository.ErrEntityNotFound) {
		require.Equalf(t, true, false, "expected error %s, actual %s", repository.ErrEntityNotFound.Error(), err.Cause().Error())
	}
}

func TestAddFlats(t *testing.T) {
	for _, flat := range flats {
		_, err := srv.CreateFlat(context.Background(), models.FlatCreateRequest{HouseID: flat.HouseID, Price: &flat.Price, Room: flat.Room})
		if err != nil {
			require.Equalf(t, true, false, "expected no error, actual %s", err.Cause())
		}
	}
}

func TestFlatsNonExistingHouse(t *testing.T) {
	_, err := srv.Flats(context.Background(), models.HouseGetFlatsRequest{ID: 2}, models.User{UserType: models.Moderator})
	if !errors.Is(err.Cause(), repository.ErrEntityNotFound) {
		require.Equalf(t, true, false, "expected error %s, actual %s", repository.ErrEntityNotFound.Error(), err.Cause().Error())
	}
}

func TestFlatsClient(t *testing.T) {
	fs, err := srv.Flats(context.Background(), models.HouseGetFlatsRequest{ID: 1}, models.User{UserType: models.Client})
	if err != nil {
		require.Equalf(t, true, false, "expected no error, actual %s", err.Cause())
	}
	require.Equal(t, 0, len(fs.Flats))
}

func TestFlatsModerator(t *testing.T) {
	fs, err := srv.Flats(context.Background(), models.HouseGetFlatsRequest{ID: 1}, models.User{UserType: models.Moderator})
	if err != nil {
		require.Equalf(t, true, false, "expected no error, actual %s", err.Cause())
	}
	require.Equal(t, flats, fs.Flats)
}
