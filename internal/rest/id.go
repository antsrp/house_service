package rest

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const requestIDKey = "request_id"

func generateRequestID(c *gin.Context) {
	go func() {
		if err := recover(); err != nil {
			slog.Info("error while generating request id", slog.Any("error", err))
		}
	}()
	id := uuid.New()
	c.Set(requestIDKey, id.String())
}
