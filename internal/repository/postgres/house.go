package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/antsrp/house_service/internal/domain/models"
	"github.com/antsrp/house_service/internal/repository"
	"github.com/jackc/pgx/v5"
)

type HouseStorage struct {
	conn Connection
}

var _ repository.HouseStorage = HouseStorage{}

func NewHouseStorage(conn Connection) HouseStorage {
	return HouseStorage{
		conn: conn,
	}
}

func (f HouseStorage) Create(ctx context.Context, req models.HouseCreateRequest) (models.HouseCreateResponse, repository.DatabaseError) {
	pattern := `INSERT INTO houses (%s) VALUES (%s) RETURNING id, created_at, updated_at`
	var (
		id                   int
		createdAt, updatedAt time.Time
	)
	fields := []string{
		"address",
		"year",
	}
	values := []string{
		fmt.Sprintf(`'%s'`, req.Address),
		strconv.Itoa(*req.Year),
	}
	if req.Developer != nil {
		fields = append(fields, "developer")
		values = append(values, fmt.Sprintf(`'%s'`, *req.Developer))
	}
	query := fmt.Sprintf(pattern, strings.Join(fields, ","), strings.Join(values, ","))
	if err := f.conn.PC.QueryRow(ctx, query).Scan(&id, &createdAt, &updatedAt); err != nil {
		return models.HouseCreateResponse{}, NewError("can't create new house", err)
	}

	house := models.House{
		ID:        id,
		Address:   req.Address,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	if req.Year != nil {
		house.Year = *req.Year
	}
	if req.Developer != nil {
		house.Developer = *req.Developer
	}

	return models.HouseCreateResponse{
		House: house,
	}, nil
}

func (f HouseStorage) Flats(ctx context.Context, req models.HouseGetFlatsRequest, user models.User) (models.HouseGetFlatsResponse, repository.DatabaseError) {
	houseQuery := `SELECT id, address, year, developer, created_at, updated_at FROM houses WHERE id = $1`
	var (
		house     models.House
		developer sql.NullString
	)
	if err := f.conn.PC.QueryRow(ctx, houseQuery, req.ID).
		Scan(&house.ID, &house.Address, &house.Year, &developer, &house.CreatedAt, &house.UpdatedAt); err != nil {
		s := fmt.Sprintf("can't get house %d", req.ID)
		if errors.Is(err, pgx.ErrNoRows) {
			return models.HouseGetFlatsResponse{}, NewError(s, repository.ErrEntityNotFound)
		}
		return models.HouseGetFlatsResponse{}, NewError(s, err)
	}
	if developer.Valid {
		house.Developer = developer.String
	}

	query := `SELECT id, house_id, price, rooms, status FROM flats WHERE house_id = $1`

	if user.UserType == models.Client {
		query = fmt.Sprintf("%s AND status = '%s'", query, models.Approved)
	}

	rows, err := f.conn.PC.Query(ctx, query, req.ID)
	if err != nil {
		return models.HouseGetFlatsResponse{}, NewError(fmt.Sprintf("can't get flats of house %d for user type %s", req.ID, user.UserType), err)
	}
	defer rows.Close()
	var flats []models.Flat
	for rows.Next() {
		var flat models.Flat
		if err := rows.Scan(&flat.ID, &flat.HouseID, &flat.Price, &flat.Room, &flat.Status); err != nil {
			return models.HouseGetFlatsResponse{}, NewError("can't scan flat", err)
		}
		flats = append(flats, flat)
	}

	return models.HouseGetFlatsResponse{
		House: house,
		Flats: flats,
	}, nil
}
