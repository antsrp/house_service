package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/antsrp/house_service/internal/domain/models"
	"github.com/antsrp/house_service/internal/repository"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type FlatStorage struct {
	conn Connection
}

var _ repository.FlatStorage = FlatStorage{}

func NewFlatStorage(conn Connection) FlatStorage {
	return FlatStorage{
		conn: conn,
	}
}

func (f FlatStorage) Create(ctx context.Context, req models.FlatCreateRequest) (models.FlatCreateResponse, repository.DatabaseError) {
	query := `INSERT INTO flats (house_id, price, rooms, status) VALUES ($1, $2, $3, $4) RETURNING id`
	var id int
	if err := f.conn.PC.QueryRow(ctx, query, req.HouseID, *(req.Price), req.Room, models.Created).Scan(&id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.ForeignKeyViolation {
			return models.FlatCreateResponse{}, NewError("can't create new flat", fmt.Errorf("house %d is not exist", req.HouseID))
		}
		return models.FlatCreateResponse{}, NewError(fmt.Sprintf("can't create new flat for house %d", req.HouseID), err)
	}

	return models.FlatCreateResponse{
		Flat: models.Flat{
			ID:      id,
			HouseID: req.HouseID,
			Price:   *req.Price,
			Room:    req.Room,
			Status:  models.Created,
		},
	}, nil
}

func (f FlatStorage) Get(ctx context.Context, req models.Flat) (models.Flat, repository.DatabaseError) {
	query := `SELECT id, house_id, price, rooms, status FROM flats WHERE id = $1`

	var result models.Flat

	if err := f.conn.PC.QueryRow(ctx, query, req.ID).Scan(&result.ID, &result.HouseID, &result.Price, &result.Room, &result.Status); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Flat{}, NewError(fmt.Sprintf("can't get flat with id %d", req.ID), repository.ErrEntityNotFound)
		}
		return models.Flat{}, NewError(fmt.Sprintf("can't get flat by id %d", req.ID), err)
	}

	return result, nil
}

func (f FlatStorage) Update(ctx context.Context, req models.FlatUpdateRequest) (models.FlatUpdateResponse, repository.DatabaseError) {
	pattern := `UPDATE flats SET %s WHERE id = $1`

	var sets []string
	if req.Price != nil {
		sets = append(sets, fmt.Sprintf("price = %d", *req.Price))
	}
	if req.Room > 0 {
		sets = append(sets, fmt.Sprintf("rooms = %d", req.Room))
	}
	if req.Status != nil {
		sets = append(sets, fmt.Sprintf(`status = '%s'`, *req.Status))
	}
	query := fmt.Sprintf(pattern, strings.Join(sets, ","))

	tag, err := f.conn.PC.Exec(ctx, query, req.ID)
	if err != nil {
		return models.FlatUpdateResponse{}, NewError("can't update flat", err)
	}

	if tag.RowsAffected() == 0 {
		return models.FlatUpdateResponse{}, NewError("can't update flat", repository.ErrNoRowsAffected)
	}

	flat, dbErr := f.Get(ctx, models.Flat{ID: req.ID})
	if dbErr != nil {
		return models.FlatUpdateResponse{}, dbErr
	}

	return models.FlatUpdateResponse{
		Flat: flat,
	}, nil
}
