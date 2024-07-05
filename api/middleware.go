package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Authenticate
func Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return next(c)
	}
}

func CustomHTTPErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	ae, ok := err.(*APIError)
	if !ok {
		ae = &APIError{
			StatusCode: http.StatusInternalServerError,
			Message:    http.StatusText(http.StatusInternalServerError),
		}
	}

	if err := c.JSON(ae.StatusCode, ae); err != nil {
		c.Logger().Error(err)
	}
}
