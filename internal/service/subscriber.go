package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/antsrp/house_service/internal/repository"
	"golang.org/x/exp/rand"
)

type SubscriberServicer interface {
	Add(ctx context.Context, email string, id int) Error
	SendEmail(ctx context.Context, id int)
}

type SubscriberService struct {
	logger  *slog.Logger
	storage repository.SubscriberStorage
}

func NewSubscriberService(logger *slog.Logger, storage repository.SubscriberStorage) SubscriberService {
	return SubscriberService{
		logger:  logger,
		storage: storage,
	}
}

var _ SubscriberServicer = SubscriberService{}

func (s SubscriberService) Add(ctx context.Context, email string, id int) Error {
	if err := s.storage.Add(ctx, email, id); err != nil {
		return NewServiceError(StatusByError(err), err.Cause(), DatabaseErrorCode)
	}
	return nil
}
func (s SubscriberService) SendEmail(ctx context.Context, id int) {
	emails, err := s.storage.Get(ctx, id)
	if err != nil {
		s.logger.Error("can't get emails of subscribers", slog.Any("error", err.Cause()))
		return
	}
	message := fmt.Sprintf("house %d is updated!", id)
	for _, email := range emails {
		if err := s.send(ctx, email, message); err != nil {
			str := fmt.Sprintf("can't send email to %s about update of house %d: %s", email, id, err.Error())
			s.logger.Error(str)
		}
	}
}

func (s SubscriberService) send(_ context.Context, recipient, message string) error {
	duration := time.Duration(rand.Int63n(3000)) * time.Millisecond
	time.Sleep(duration)

	// Имитация неуспешной отправки сообщения
	errorProbability := 0.1
	if rand.Float64() < errorProbability {
		return fmt.Errorf("internal error")
	}

	s.logger.Info(fmt.Sprintf("send message '%s' to '%s'", message, recipient))

	return nil
}
