package service

import (
	"fmt"

	"github.com/antsrp/house_service/internal/repository"
)

type ErrorStatus = string

const (
	BadRequest ErrorStatus = "bad request"
	Conflict   ErrorStatus = "conflict"
	Internal   ErrorStatus = "internal error"
)

type ErrorCode = int

const (
	DatabaseErrorCode ErrorCode = iota + 1
	CreateTokenErrorCode
	CryptoErrorCode
)

type Error interface {
	Status() ErrorStatus
	Cause() error
	Code() ErrorCode
}

func StatusByError(err repository.DatabaseError) ErrorStatus {
	status := BadRequest
	if err.IsInternal() {
		status = Internal
	}
	return status
}

type serviceError struct {
	status ErrorStatus
	real   error
	code   ErrorCode
}

//var defaultInternalError = serviceError{status: Internal, real: ErrDefaultInternalError}

func (e serviceError) Status() ErrorStatus {
	return e.status
}

func (e serviceError) Cause() error {
	return e.real
}

func (e serviceError) Code() ErrorCode {
	return e.code
}

func NewServiceError(status ErrorStatus, err error, code ErrorCode) Error {
	return serviceError{
		status: status,
		real:   err,
		code:   code,
	}
}

var (
	ErrDefaultInternalError = fmt.Errorf("internal server error, try again later")
	ErrHouseNotFound        = fmt.Errorf("house not found")
	ErrFlatNotFound         = fmt.Errorf("flat not found")
	ErrOnModeration         = fmt.Errorf("flat is already on moderation")
	ErrUserAlreadyExists    = fmt.Errorf("user already exists")
	ErrUserNotFound         = fmt.Errorf("user not found")
)
