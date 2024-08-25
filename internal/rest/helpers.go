package rest

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/antsrp/house_service/internal/domain/models"
	"github.com/antsrp/house_service/internal/service"
	"github.com/gin-gonic/gin"
)

func abort(c *gin.Context, ctx context.Context, logger *slog.Logger, level slog.Level, msg string, extra map[string]any, status int, codes ...service.ErrorCode) {
	var f func(context.Context, string, ...any)
	switch level {
	case slog.LevelDebug:
		f = logger.DebugContext
	case slog.LevelError:
		f = logger.ErrorContext
	default:
		f = logger.InfoContext
	}

	var attrs []any
	for k, v := range extra {
		attrs = append(attrs, slog.Any(k, v))
	}
	f(ctx, msg, attrs...)
	if status == http.StatusInternalServerError {
		requestID, _ := c.Get(requestIDKey)
		m := gin.H{"message": service.ErrDefaultInternalError.Error(), "request_id": requestID}
		if len(codes) != 0 {
			m["code"] = codes[0]
		}
		c.AbortWithStatusJSON(status, m)
		return
	}
	c.AbortWithStatus(status)
}

func parseRequestContext(c *gin.Context, logger *slog.Logger) context.Context {
	val, ok := c.Get(requestContextKey)
	if !ok {
		logger.Info("no context attached to current request")
		return nil
	}
	ctx, ok := val.(context.Context)
	if !ok {
		logger.Info("no context attached to current request")
		return nil
	}
	return ctx
}

func codeByStatus(status service.ErrorStatus) int {
	var code int
	switch status {
	case service.Conflict:
		code = http.StatusConflict
	case service.Internal:
		code = http.StatusInternalServerError
	case service.BadRequest:
		code = http.StatusBadRequest
	}

	return code
}

var (
	errNotProvided  = fmt.Errorf("parameter is not provided")
	errNotAnInteger = fmt.Errorf("parameter is not an integer")
)

func paramInt(c *gin.Context, name string) (int, error) {
	s := c.Param(name)
	if s == "" {
		return 0, errNotProvided
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0, errNotAnInteger
	}

	return val, nil
}

func paramIntErrorHandler(c *gin.Context, ctx context.Context, logger *slog.Logger, err error, field string) {
	if errors.Is(err, errNotProvided) {
		abort(c, ctx, logger, slog.LevelInfo, fmt.Sprintf("%s is not provided", field), nil, http.StatusBadRequest)
	} else if errors.Is(err, errNotAnInteger) {
		abort(c, ctx, logger, slog.LevelInfo, fmt.Sprintf("%s is not an integer", field), nil, http.StatusBadRequest)
	} else { // some another error
		abort(c, ctx, logger, slog.LevelError, fmt.Sprintf("cannot parse id parameter: %s", err.Error()), nil, http.StatusInternalServerError)
	}
}

var (
	errTagNotExist       = fmt.Errorf("user tag not exist")
	errTagDataNoUserType = fmt.Errorf("user tag data is not user type")
)

const authusertag = "auth-user-tag-data"

func setUser(c *gin.Context, user models.User) {
	c.Set(authusertag, user)
}

func parseUser(c *gin.Context) (models.User, error) {
	data, exists := c.Get(authusertag)
	if !exists {
		return models.User{}, errTagNotExist
	}
	user, ok := data.(models.User)
	if !ok {
		return models.User{}, errTagDataNoUserType
	}
	return user, nil
}

func parseUserErrorHandler(c *gin.Context, ctx context.Context, logger *slog.Logger, err error) {
	if errors.Is(err, errNotProvided) {
		abort(c, ctx, logger, slog.LevelInfo, "id of house is not provided", nil, http.StatusBadRequest)
	} else if errors.Is(err, errNotAnInteger) {
		abort(c, ctx, logger, slog.LevelInfo, "id of house is not an integer", nil, http.StatusBadRequest)
	} else { // some another error
		abort(c, ctx, logger, slog.LevelError, fmt.Sprintf("cannot parse id parameter: %s", err.Error()), nil, http.StatusInternalServerError)
	}
}
