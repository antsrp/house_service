package repository

import "fmt"

const (
	msgNoRowsAffected = "no rows affected"
	msgEntityNotFound = "no entity found"
	msgAlreadyExists  = "entity already exists"
)

var (
	ErrNoRowsAffected      = fmt.Errorf(msgNoRowsAffected)
	ErrEntityNotFound      = fmt.Errorf(msgEntityNotFound)
	ErrEntityAlreadyExists = fmt.Errorf(msgAlreadyExists)
)

type DatabaseError interface {
	IsInternal() bool
	Cause() error
}
