package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/antsrp/house_service/internal/domain/models"
	"github.com/antsrp/house_service/internal/repository"
	"github.com/jackc/pgx/v5"
)

type TokenStorage struct {
	conn Connection
}

var _ repository.TokenStorage = TokenStorage{}

func NewTokenStorage(conn Connection) TokenStorage {
	return TokenStorage{
		conn: conn,
	}
}

func (t TokenStorage) UserByToken(ctx context.Context, token string) (models.User, repository.DatabaseError) {
	query := `SELECT users.id, user_type FROM tokens 
	JOIN users ON users.id = tokens.user_id
	WHERE token = $1 AND created_at BETWEEN NOW() - INTERVAL '24 HOURS' AND NOW()`

	var user models.User

	if err := t.conn.PC.QueryRow(ctx, query, token).Scan(&user.ID, &user.UserType); err != nil {
		s := fmt.Sprintf("can't get user by token %s", token)
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, NewError(s, repository.ErrEntityNotFound)
		}
		return models.User{}, NewError(s, err)
	}

	return user, nil
}
func (t TokenStorage) TokenByID(ctx context.Context, user models.User) (string, repository.DatabaseError) {
	query := `SELECT token FROM tokens 
	JOIN users ON users.id = tokens.user_id
	WHERE users.id = $1`

	var token string

	if err := t.conn.PC.QueryRow(ctx, query, user.ID).Scan(&token); err != nil {
		s := fmt.Sprintf("can't get token for user by id %v", user.ID)
		if errors.Is(err, pgx.ErrNoRows) {
			return "", NewError(s, repository.ErrEntityNotFound)
		}
		return "", NewError(s, err)
	}

	return token, nil
}
func (t TokenStorage) AddToken(ctx context.Context, user models.User, token string, currentTime time.Time) repository.DatabaseError {
	query := `INSERT INTO tokens (user_id, token, created_at) VALUES ($1, $2, $3)
	ON CONFLICT (user_id) DO UPDATE
	SET token = $2, created_at = $3`

	if _, err := t.conn.PC.Exec(ctx, query, user.ID, token, currentTime); err != nil {
		return NewError(fmt.Sprintf("can't insert or update token for user %v", user.ID), err)
	}

	return nil
}
