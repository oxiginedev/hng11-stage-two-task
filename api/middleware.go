package api

import (
	"context"
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

		_, err = a.database.FetchUserByID(c.Request().Context(), vt.UserID)
		if err != nil {
			return newAPIError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		}

		// authCtx := context.WithValue(c.Request().Context(), AuthUserCtx, user)

		return next(c)
	}
}

func setAuthUserInContext(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, AuthUserCtx, user)
}

func GetAuthUserFromContext(ctx context.Context) *models.User {
	return ctx.Value(AuthUserCtx).(*models.User)
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
