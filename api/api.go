package api

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/oxiginedev/hng11-stage-two-task/config"
	"github.com/oxiginedev/hng11-stage-two-task/pkg/datastore"
)

type API struct {
	database datastore.Datastore
	config   *config.Configuration
}

type APIError struct {
	StatusCode int
	Message    string
}

func (a *APIError) Error() string {
	return fmt.Sprintf("api error | statusCode: %d | message: %s", a.StatusCode, a.Message)
}

// New instantiates API
func New(db datastore.Datastore, cfg *config.Configuration) *API {
	return &API{
		database: db,
		config:   cfg,
	}
}

func (a *API) Routes() *echo.Echo {
	r := echo.New()

	r.Use(middleware.Logger())
	r.Use(middleware.Recover())

	r.POST("/auth/register", a.HandleRegister)
	r.POST("/auth/login", a.HandleLogin)

	r.HTTPErrorHandler = CustomHTTPErrorHandler

	return r
}
