package api

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/oxiginedev/hng11-stage-two-task/config"
	"github.com/oxiginedev/hng11-stage-two-task/pkg/datastore"
	"github.com/oxiginedev/hng11-stage-two-task/types"
)

type (
	API struct {
		database datastore.Datastore
		config   *config.Configuration
	}

	APIError struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Code    int    `json:"statusCode"`
	}

	APIValidator struct {
		validator *validator.Validate
	}
)

func (av *APIValidator) Validate(i interface{}) error {
	err := av.validator.Struct(i)
	if err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, ve.Error())
		}
		return err
	}
	return nil
}

func (a *APIError) Error() string {
	return fmt.Sprintf("api error | statusCode: %d | message: %s", a.Code, a.Message)
}

func newAPIError(code int, message string) *APIError {
	return &APIError{
		Status:  http.StatusText(code),
		Message: message,
		Code:    code,
	}
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

	r.HTTPErrorHandler = CustomHTTPErrorHandler
	r.Validator = &APIValidator{validator: validator.New()}

	r.Use(middleware.Logger())
	r.Use(middleware.Recover())

	r.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, types.M{
			"status":  "success",
			"message": "HNG11 Stage 2 Task",
			"data": types.M{
				"slackId": "0xigine",
				"email":   "adedaramolaadetimehin@gmail.com",
			},
		})
	})

	r.POST("/auth/register", a.HandleRegister)
	r.POST("/auth/login", a.HandleLogin)

	r.GET("/api/users/:id", a.HandleGetUser, a.Authenticate)

	r.GET("/api/organisations", a.HandleGetOrganisations, a.Authenticate)
	r.GET("/api/organisations/:orgId", a.HandleGetOrganisation, a.Authenticate)
	r.POST("/api/organisations/:orgId/users", a.HandleCreateOrganisationUser, a.Authenticate)
	r.POST("/api/organisations", a.HandleCreateOrganisation, a.Authenticate)

	return r
}
