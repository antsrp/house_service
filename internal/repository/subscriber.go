package repository

import "context"

type SubscriberStorage interface {
	Add(ctx context.Context, email string, id int) DatabaseError
	Get(ctx context.Context, id int) ([]string, DatabaseError) // emails slice
}
