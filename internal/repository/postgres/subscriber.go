package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/antsrp/house_service/internal/repository"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type SubscriberStorage struct {
	conn Connection
}

var _ repository.SubscriberStorage = SubscriberStorage{}

func NewSubscriberStorage(conn Connection) SubscriberStorage {
	return SubscriberStorage{
		conn: conn,
	}
}

func (s SubscriberStorage) Add(ctx context.Context, email string, id int) repository.DatabaseError {
	query := `INSERT INTO subscribers (email, house_id) VALUES ($1, $2)`

	if _, err := s.conn.PC.Exec(ctx, query, email, id); err != nil {
		s := fmt.Sprintf("can't add new subscriber with email %s to house %d", email, id)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return NewError(s, repository.ErrEntityAlreadyExists)
		}
		return NewError(s, err)
	}

	return nil
}

func (s SubscriberStorage) Get(ctx context.Context, id int) ([]string, repository.DatabaseError) {
	query := `SELECT email FROM subscribers WHERE house_id = $1`

	rows, err := s.conn.PC.Query(ctx, query, id)
	if err != nil {
		return nil, NewError(fmt.Sprintf("can't get subscribers of house %d", id), err)
	}
	defer rows.Close()
	var emails []string
	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			return nil, NewError(fmt.Sprintf("can't get subscriber of house %d", id), err)
		}
		emails = append(emails, email)
	}

	return emails, nil
}
