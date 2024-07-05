package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/oxiginedev/hng11-stage-two-task/types"
)

func (a *API) HandleGetUser(c echo.Context) error {
	user := GetAuthUserFromContext(c.Request().Context())

	data := types.M{
		"status":  "success",
		"message": "Profile retrieved",
		"data":    user,
	}

	return c.JSON(http.StatusOK, data)
}
