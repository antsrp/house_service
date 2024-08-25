package rest

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/antsrp/house_service/internal/domain/models"
	"github.com/antsrp/house_service/internal/service"
	rs "github.com/antsrp/house_service/pkg/infrastructure/rest"
	"github.com/antsrp/house_service/pkg/log"
	"github.com/gin-gonic/gin"
)

const requestContextKey string = "req-context-key"

type Handler struct {
	engine           *gin.Engine
	settings         rs.Settings
	logger           *slog.Logger
	authHandler      authHandler
	houseFlatService service.HouseFlatServicer
	userService      service.UserServicer
}

func NewHandler(logger *slog.Logger, settings rs.Settings, hfService service.HouseFlatServicer, userService service.UserServicer, tokenService service.TokenServicer) Handler {
	h := Handler{
		logger:           logger,
		engine:           gin.Default(),
		settings:         settings,
		authHandler:      newAuthHandler(logger, tokenService),
		houseFlatService: hfService,
		userService:      userService,
	}
	h.routes()

	return h
}

func (h Handler) provideRequestID(c *gin.Context) {
	val, found := c.Get(requestIDKey)
	if !found {
		h.logger.Info("request id is not provided")
		return
	}
	ctx := log.AppendCtx(context.Background(), slog.Any(requestIDKey, val))
	c.Set(requestContextKey, ctx)

	h.logger.InfoContext(ctx, fmt.Sprintf("new request registered, path %s", c.Request.URL.Path))
}

func (h Handler) routes() {
	group := h.engine.Group("/", generateRequestID, h.provideRequestID)
	group.GET("/dummyLogin", h.dummyLogin)
	group.POST("/login", h.login)
	group.POST("/register", h.register)
	houseGroup, flatGroup := group.Group("/house", h.authHandler.authRequired), group.Group("/flat", h.authHandler.authRequired)
	houseGroup.POST("/create", h.authHandler.moderatorAuthRequired, h.houseCreate)
	houseGroup.GET("/:id", h.houseByID)
	houseGroup.POST("/:id/subscribe", h.subscribe)
	flatGroup.POST("/create", h.flatCreate)
	flatGroup.POST("/update", h.authHandler.moderatorAuthRequired, h.flatUpdate)
}

func (h Handler) Run() error {
	address := fmt.Sprintf("%s:%s", h.settings.Host, h.settings.Port)
	if err := h.engine.Run(address); err != nil {
		return fmt.Errorf("can't run server: %w", err)
	}
	return nil
}

func (h Handler) dummyLogin(c *gin.Context) { // GET /dummyLogin
	ctx := parseRequestContext(c, h.logger)
	ut, found := c.GetQuery("user_type")
	if !found {
		abort(c, ctx, h.logger, slog.LevelInfo, "user type is not provided", nil, http.StatusBadRequest)
		return
	}

	var token string
	if ut == models.Client {
		token = h.authHandler.dummyClient()
	} else if ut == models.Moderator {
		token = h.authHandler.dummyModerator()
	}
	if token == "" {
		abort(c, ctx, h.logger, slog.LevelInfo, "invalid user type is provided", nil, http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, models.DummyLoginResponse{Token: token})
}

func (h Handler) login(c *gin.Context) { // POST /login
	ctx := parseRequestContext(c, h.logger)

	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		abort(c, ctx, h.logger, slog.LevelInfo, "cannot parse request data", map[string]any{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	if req.UserID == nil {
		abort(c, ctx, h.logger, slog.LevelInfo, "user id is not provided", nil, http.StatusBadRequest)
		return
	}

	if req.Password == nil {
		abort(c, ctx, h.logger, slog.LevelInfo, "password is not provided", nil, http.StatusBadRequest)
		return
	}

	resp, err := h.userService.Login(ctx, req)
	if err != nil {
		abort(c, ctx, h.logger, slog.LevelInfo, "cannot login user", map[string]any{"error": err.Cause()}, codeByStatus(err.Status()), err.Code())
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h Handler) register(c *gin.Context) { // POST /register
	ctx := parseRequestContext(c, h.logger)

	var req models.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		abort(c, ctx, h.logger, slog.LevelInfo, "cannot parse request data", map[string]any{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	if req.Email == nil {
		abort(c, ctx, h.logger, slog.LevelInfo, "email is not provided", nil, http.StatusBadRequest)
		return
	}

	if req.Password == nil {
		abort(c, ctx, h.logger, slog.LevelInfo, "password is not provided", nil, http.StatusBadRequest)
		return
	}

	if req.UserType == nil {
		abort(c, ctx, h.logger, slog.LevelInfo, "user type is not provided", nil, http.StatusBadRequest)
		return
	}

	resp, err := h.userService.Register(ctx, req)
	if err != nil {
		abort(c, ctx, h.logger, slog.LevelInfo, "cannot register user", map[string]any{"error": err.Cause()}, codeByStatus(err.Status()), err.Code())
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h Handler) houseCreate(c *gin.Context) { // POST /house/create
	ctx := parseRequestContext(c, h.logger)
	var req models.HouseCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		abort(c, ctx, h.logger, slog.LevelInfo, "cannot parse request data", map[string]any{"error": err.Error()}, http.StatusBadRequest)
		return
	}
	if req.Address == "" {
		abort(c, ctx, h.logger, slog.LevelInfo, "address for house is not provided", nil, http.StatusBadRequest)
		return
	}
	if req.Year == nil {
		abort(c, ctx, h.logger, slog.LevelInfo, "year of creation is not provided", nil, http.StatusBadRequest)
		return
	}
	if *(req.Year) < 0 {
		abort(c, ctx, h.logger, slog.LevelInfo, "year value is unacceptable", nil, http.StatusBadRequest)
		return
	}

	house, err := h.houseFlatService.CreateHouse(ctx, req)
	if err != nil {
		abort(c, ctx, h.logger, slog.LevelInfo, err.Cause().Error(), nil, codeByStatus(err.Status()), err.Code())
		return
	}

	c.JSON(http.StatusOK, house)
}

func (h Handler) houseByID(c *gin.Context) { // GET /house/{id}
	ctx := parseRequestContext(c, h.logger)
	id, err := paramInt(c, "id")
	if err != nil {
		paramIntErrorHandler(c, ctx, h.logger, err, "id of house")
		return
	}
	user, err := parseUser(c)
	if err != nil {
		parseUserErrorHandler(c, ctx, h.logger, err)
		return
	}

	flats, srvErr := h.houseFlatService.Flats(ctx, models.HouseGetFlatsRequest{ID: id}, user)
	if srvErr != nil {
		abort(c, ctx, h.logger, slog.LevelInfo, srvErr.Cause().Error(), nil, codeByStatus(srvErr.Status()), srvErr.Code())
		return
	}

	c.JSON(http.StatusOK, flats)
}

func (h Handler) subscribe(c *gin.Context) { // POST /house/{id}/subscribe
	ctx := parseRequestContext(c, h.logger)
	id, err := paramInt(c, "id")
	if err != nil {
		paramIntErrorHandler(c, ctx, h.logger, err, "id of house")
		return
	}

	var req models.SubscribeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		abort(c, ctx, h.logger, slog.LevelInfo, "cannot parse request data", map[string]any{"error": err.Error()}, http.StatusBadRequest)
		return
	}
	if req.Email == nil {
		abort(c, ctx, h.logger, slog.LevelInfo, "email is not provided", nil, http.StatusBadRequest)
		return
	}
	if err := h.houseFlatService.AddSubscriber(ctx, *req.Email, id); err != nil {
		abort(c, ctx, h.logger, slog.LevelInfo, err.Cause().Error(), nil, codeByStatus(err.Status()), err.Code())
		return
	}
	c.Status(http.StatusOK)
}

func (h Handler) flatCreate(c *gin.Context) { // POST /flat/create
	ctx := parseRequestContext(c, h.logger)

	var req models.FlatCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		abort(c, ctx, h.logger, slog.LevelInfo, "cannot parse request data", map[string]any{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	if req.HouseID < 1 {
		abort(c, ctx, h.logger, slog.LevelInfo, "house_id field is not set or value is inappropriate", nil, http.StatusBadRequest)
		return
	}

	if req.Price == nil || *req.Price < 0 {
		abort(c, ctx, h.logger, slog.LevelInfo, "price field is not set or value is inappropriate", nil, http.StatusBadRequest)
		return
	}

	if req.Room < 1 {
		abort(c, ctx, h.logger, slog.LevelInfo, "room field is not set or value is inappropriate", nil, http.StatusBadRequest)
		return
	}

	flat, err := h.houseFlatService.CreateFlat(ctx, req)
	if err != nil {
		abort(c, ctx, h.logger, slog.LevelInfo, err.Cause().Error(), nil, codeByStatus(err.Status()), err.Code())
		return
	}

	c.JSON(http.StatusOK, flat)
}

func (h Handler) flatUpdate(c *gin.Context) { // POST /flat/update
	ctx := parseRequestContext(c, h.logger)

	var req models.FlatUpdateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		abort(c, ctx, h.logger, slog.LevelInfo, "cannot parse request data", map[string]any{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	if req.ID < 1 {
		abort(c, ctx, h.logger, slog.LevelInfo, "id field is not set or value is inappropriate", nil, http.StatusBadRequest)
		return
	}

	if req.Price == nil || *req.Price < 0 {
		abort(c, ctx, h.logger, slog.LevelInfo, "price field is not set or value is inappropriate", nil, http.StatusBadRequest)
		return
	}

	if req.Room < 1 {
		abort(c, ctx, h.logger, slog.LevelInfo, "room field is not set or value is inappropriate", nil, http.StatusBadRequest)
		return
	}

	flat, err := h.houseFlatService.UpdateFlat(ctx, req)
	if err != nil {
		abort(c, ctx, h.logger, slog.LevelInfo, err.Cause().Error(), nil, codeByStatus(err.Status()), err.Code())
		return
	}

	c.JSON(http.StatusOK, flat)
}
