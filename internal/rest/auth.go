package rest

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/antsrp/house_service/internal/domain/models"
	"github.com/antsrp/house_service/internal/service"
	"github.com/gin-gonic/gin"
)

const (
	defaultClientToken    = "default-client-token"
	defaultModeratorToken = "default-moderator-token"
)

type authHandler struct {
	logger       *slog.Logger
	tokenService service.TokenServicer
}

func newAuthHandler(logger *slog.Logger, tokenService service.TokenServicer) authHandler {
	return authHandler{
		logger:       logger,
		tokenService: tokenService,
	}
}

var errInternal = fmt.Errorf("internal auth error")

func (h authHandler) parseToken(c *gin.Context) error {
	header := c.GetHeader("Authorization")
	if header == "" {
		return fmt.Errorf("no token provided")
	}

	parts := strings.Split(header, " ")
	if len(parts) != 2 {
		return fmt.Errorf("invalid token data")
	}

	tp, token := parts[0], parts[1]

	if tp != "Bearer" {
		return fmt.Errorf("invalid token data: not a bearer token")
	}

	if h.checkDummy(c, token) {
		return nil
	}
	ctx := context.Background()

	user, err := h.tokenService.UserByToken(ctx, token)
	if err != nil {
		if err.Status() == service.Internal {
			abort(c, ctx, h.logger, slog.LevelInfo, "cannot recognize user by token", map[string]any{"error": err.Cause()}, codeByStatus(err.Status()), err.Code())
			return errInternal
		}
		return err.Cause()
	}

	setUser(c, user)
	return nil

	//return fmt.Errorf("only dummy users available for now")
}

func (h authHandler) dummyClient() string {
	return defaultClientToken
}

func (h authHandler) dummyModerator() string {
	return defaultModeratorToken
}

func (h authHandler) checkDummy(c *gin.Context, token string) bool {
	var user models.User
	if token == defaultClientToken {
		user.UserType = models.Client
	} else if token == defaultModeratorToken {
		user.UserType = models.Moderator
	} else {
		return false
	}
	setUser(c, user)
	return true
}

func (h authHandler) authRequired(c *gin.Context) {
	ctx := parseRequestContext(c, h.logger)
	if err := h.parseToken(c); err != nil {
		if !errors.Is(err, errInternal) {
			abort(c, ctx, h.logger, slog.LevelInfo, "user is not authorized", map[string]any{"error": err.Error()}, http.StatusUnauthorized)
		}
		return
	}
}

func (h authHandler) moderatorAuthRequired(c *gin.Context) {
	ctx := parseRequestContext(c, h.logger)
	user, err := parseUser(c)
	if err != nil {
		parseUserErrorHandler(c, ctx, h.logger, err)
		return
	}
	if user.UserType != models.Moderator {
		abort(c, ctx, h.logger, slog.LevelInfo, "current user is not an admin", nil, http.StatusUnauthorized)
		return
	}
	//c.Set(authusertag, user)
}
