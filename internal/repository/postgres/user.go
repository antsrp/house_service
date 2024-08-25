package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/antsrp/house_service/internal/domain/models"
	"github.com/antsrp/house_service/internal/repository"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type UserStorage struct {
	conn Connection
}

var _ repository.UserStorage = UserStorage{}

func NewUserStorage(conn Connection) UserStorage {
	return UserStorage{
		conn: conn,
	}
}

func (u UserStorage) Add(ctx context.Context, user models.User) (models.User, repository.DatabaseError) {
	query := `INSERT INTO users (email, password, user_type) VALUES ($1, $2, $3) RETURNING id`

	var result models.User

	if err := u.conn.PC.QueryRow(ctx, query, user.Email, user.Password, user.UserType).Scan(&result.ID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return models.User{}, NewError(fmt.Sprintf("can't add new user with email %s", user.Email), repository.ErrEntityAlreadyExists)
		}
		return models.User{}, NewError("can't add new user", err)
	}

	return result, nil
}
func (u UserStorage) Get(ctx context.Context, user models.User) (models.User, repository.DatabaseError) {
	query := `SELECT id, email, user_type FROM users WHERE id = $1 AND password = $2`

	var result models.User

	if err := u.conn.PC.QueryRow(ctx, query, user.ID, user.Password).Scan(&result.ID, &result.Email, &result.UserType); err != nil {
		s := fmt.Sprintf("can't get user with id %v", user.ID)
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, NewError(s, repository.ErrEntityNotFound)
		}
		return models.User{}, NewError(s, err)
	}

	return result, nil
}
