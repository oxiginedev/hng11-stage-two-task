package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/oxiginedev/hng11-stage-two-task/pkg/jwt"
	"github.com/oxiginedev/hng11-stage-two-task/pkg/models"
)

type ContextKey string

const AuthUserCtx ContextKey = "auth:user"

// Authenticate
func (a *API) Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		h := c.Request().Header.Get("Authorization")
		hStruct := strings.Split(h, " ")

		if len(hStruct) != 2 {
			return newAPIError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		}

		accessToken := hStruct[1]
		if len(strings.TrimSpace(accessToken)) == 0 {
			return newAPIError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		}

		// validate retrieved access token
		vt, err := jwt.ValidateAccessToken(a.config.JWT.Secret, accessToken)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				return newAPIError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			}
			return err
		}

		user, err := a.database.FetchUserByID(c.Request().Context(), vt.UserID)
		if err != nil {
			return newAPIError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		}

		c.Set(string(AuthUserCtx), user)

		return next(c)
	}
}

func GetAuthUserFromContext(c echo.Context) *models.User {
	return c.Get(string(AuthUserCtx)).(*models.User)
}

func CustomHTTPErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	ae, ok := err.(*APIError)
	if !ok {
		ae = newAPIError(http.StatusInternalServerError, http.StatusText(500))
	}

	if err := c.JSON(ae.Code, ae); err != nil {
		c.Logger().Error(err)
	}
}
