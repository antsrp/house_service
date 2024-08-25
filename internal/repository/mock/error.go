package mock

import "github.com/antsrp/house_service/internal/repository"

type MockError struct {
	isInternal bool
	real       error
}

func NewMockError(isInternal bool, err error) MockError {
	return MockError{
		isInternal: isInternal,
		real:       err,
	}
}

func (m MockError) IsInternal() bool {
	return m.isInternal
}

func (m MockError) Cause() error {
	return m.real
}

var _ repository.DatabaseError = MockError{}
