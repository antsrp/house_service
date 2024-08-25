package service

import (
	"context"
)

type MockSubscriberServicer interface {
	Add(ctx context.Context, email string, id int) Error
	SendEmail(ctx context.Context, id int)
}

type MockSubscriberService struct {
}

func NewMockSubscriberService() MockSubscriberService {
	return MockSubscriberService{}
}

var _ SubscriberServicer = MockSubscriberService{}

func (s MockSubscriberService) Add(ctx context.Context, email string, id int) Error {
	return nil
}
func (s MockSubscriberService) SendEmail(ctx context.Context, id int) {

}
