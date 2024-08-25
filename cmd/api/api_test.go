package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/antsrp/house_service/internal/domain/models"
	"github.com/antsrp/house_service/internal/setup"
	"github.com/stretchr/testify/require"
)

var (
	addr string

	houses []models.House
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../..") // get back to root location
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	//h, dbConnection, logger, err := setup.Setup("TEST_DB")
	h, srvSettings, dbConnection, logger, err := setup.Setup("DB")
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	defer dbConnection.Close()

	addr = fmt.Sprintf("%s://%s:%s", "http", srvSettings.Host, srvSettings.Port)

	go func() {
		if err := h.Run(); err != nil {
			logger.Error("cannot run rest server", slog.Any("error", err.Error()))
		}
	}()

	m.Run()
}

func TestAddHouseNoAuth(t *testing.T) {
	year := 2023
	developer := `PEEK`
	input := models.HouseCreateRequest{Address: "addr1", Year: &year, Developer: &developer}
	data, _ := json.Marshal(input)
	req, _ := http.NewRequest(http.MethodPost, addr+"/house/create", bytes.NewBuffer(data))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	require.Equalf(t, http.StatusUnauthorized, resp.StatusCode, "expected status code %d, actual %d", http.StatusUnauthorized, resp.StatusCode)
}

const (
	clientToken    = "default-client-token"
	moderatorToken = "default-moderator-token"
)

func TestAddHouse(t *testing.T) {
	year := 2023
	developer := `PEEK`
	input := models.HouseCreateRequest{Address: "addr1", Year: &year, Developer: &developer}
	data, _ := json.Marshal(input)
	req, _ := http.NewRequest(http.MethodPost, addr+"/house/create", bytes.NewBuffer(data))
	req.Header.Set("Authorization", "Bearer "+moderatorToken)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equalf(t, http.StatusOK, resp.StatusCode, "expected status code %d, actual %d", http.StatusOK, resp.StatusCode)

	arr, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var output models.HouseCreateResponse
	err = json.Unmarshal(arr, &output)
	require.NoError(t, err)

	houses = append(houses, output.House)
}

func TestAddFlatNonExistingHouse(t *testing.T) {
	if len(houses) == 0 {
		return
	}
	price := 1323
	input := models.FlatCreateRequest{HouseID: houses[0].ID + 15, Price: &price, Room: 2}
	data, _ := json.Marshal(input)
	req, _ := http.NewRequest(http.MethodPost, addr+"/flat/create", bytes.NewBuffer(data))
	req.Header.Set("Authorization", "Bearer "+clientToken)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equalf(t, http.StatusBadRequest, resp.StatusCode, "expected status code %d, actual %d", http.StatusBadRequest, resp.StatusCode)
}

func TestAddFlat(t *testing.T) {
	price := 1323
	input := models.FlatCreateRequest{HouseID: houses[0].ID, Price: &price, Room: 2}
	data, _ := json.Marshal(input)
	req, _ := http.NewRequest(http.MethodPost, addr+"/flat/create", bytes.NewBuffer(data))
	req.Header.Set("Authorization", "Bearer "+clientToken)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equalf(t, http.StatusOK, resp.StatusCode, "expected status code %d, actual %d", http.StatusOK, resp.StatusCode)
}

func TestFlatsNonExistingHouse(t *testing.T) {
	input := models.HouseGetFlatsRequest{ID: houses[0].ID + 222}
	data, _ := json.Marshal(input)
	req, _ := http.NewRequest(http.MethodPost, addr+"/flat/create", bytes.NewBuffer(data))
	req.Header.Set("Authorization", "Bearer "+moderatorToken)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equalf(t, http.StatusBadRequest, resp.StatusCode, "expected status code %d, actual %d", http.StatusBadRequest, resp.StatusCode)
}

func TestFlatsClient(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf(addr+"/house/%d", houses[0].ID), nil)
	req.Header.Set("Authorization", "Bearer "+clientToken)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equalf(t, http.StatusOK, resp.StatusCode, "expected status code %d, actual %d", http.StatusOK, resp.StatusCode)

	arr, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var output models.HouseGetFlatsResponse
	err = json.Unmarshal(arr, &output)
	require.NoError(t, err)

	require.Equal(t, 0, len(output.Flats))
}

func TestFlatsModerator(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf(addr+"/house/%d", houses[0].ID), nil)
	req.Header.Set("Authorization", "Bearer "+moderatorToken)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equalf(t, http.StatusOK, resp.StatusCode, "expected status code %d, actual %d", http.StatusOK, resp.StatusCode)

	arr, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var output models.HouseGetFlatsResponse
	err = json.Unmarshal(arr, &output)
	require.NoError(t, err)

	require.Equal(t, 1, len(output.Flats))
}
