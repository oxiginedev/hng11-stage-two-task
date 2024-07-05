package api

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/oklog/ulid/v2"
	"github.com/oxiginedev/hng11-stage-two-task/pkg/datastore"
	"github.com/oxiginedev/hng11-stage-two-task/pkg/models"
	"github.com/oxiginedev/hng11-stage-two-task/types"
)

type CreateOrgUserPayload struct {
	UserID string `json:"userId" validate:"required"`
}

func (a *API) HandleCreateOrganisationUser(c echo.Context) error {
	var p CreateOrgUserPayload

	if err := c.Bind(&p); err != nil {
		return err
	}

	if err := c.Validate(p); err != nil {
		return err
	}

	user := GetAuthUserFromContext(c.Request().Context())
	orgID := c.Param("orgId")

	_, err := a.database.FetchOrganisationUserByUserID(
		c.Request().Context(),
		user.ID,
		orgID,
	)
	if err != nil {
		if errors.Is(err, datastore.ErrRecordNotFound) {
			// returning 404, this seems better for security
			return newAPIError(http.StatusNotFound, "Organisation not found")
		}
		return err
	}

	ou := &models.OrganisationUser{
		ID:             ulid.Make().String(),
		OrganisationID: orgID,
		UserID:         p.UserID,
	}

	err = a.database.CreateOrganisationUser(c.Request().Context(), ou)
	if err != nil {
		return err
	}

	data := types.M{
		"status":  "success",
		"message": "User added to organisation successfully",
	}

	return c.JSON(http.StatusOK, data)
}
