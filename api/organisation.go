package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/oklog/ulid/v2"
	"github.com/oxiginedev/hng11-stage-two-task/pkg/datastore"
	"github.com/oxiginedev/hng11-stage-two-task/pkg/models"
	"github.com/oxiginedev/hng11-stage-two-task/types"
)

func (a *API) HandleGetOrganisations(c echo.Context) error {
	user := GetAuthUserFromContext(c)

	orgs, err := a.database.FetchUserOrganisations(c.Request().Context(), user.ID)
	if err != nil {
		return err
	}

	data := types.M{
		"status":  "success",
		"message": "Organisations retrieved",
		"data": types.M{
			"organisations": orgs,
		},
	}

	return c.JSON(http.StatusOK, data)
}

func (a *API) HandleGetOrganisation(c echo.Context) error {
	user := GetAuthUserFromContext(c)
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

	org, err := a.database.FetchOrganisationByID(c.Request().Context(), orgID)
	if err != nil {
		if errors.Is(err, datastore.ErrRecordNotFound) {
			return newAPIError(http.StatusNotFound, "Organisation not found")
		}
		return err
	}

	data := types.M{
		"status":  "success",
		"message": "Organisation retrieved",
		"data":    org,
	}

	return c.JSON(http.StatusOK, data)
}

type CreateOrgPayload struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

func (a *API) HandleCreateOrganisation(c echo.Context) error {
	var p CreateOrgPayload

	user := GetAuthUserFromContext(c)

	err := c.Bind(&p)
	if err != nil {
		return err
	}

	if err := c.Validate(p); err != nil {
		return err
	}

	org := &models.Organisation{
		ID:          ulid.Make().String(),
		Name:        p.Name,
		Description: p.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = a.database.CreateOrganisation(c.Request().Context(), org)
	if err != nil {
		return err
	}

	oru := &models.OrganisationUser{
		ID:             ulid.Make().String(),
		OrganisationID: org.ID,
		UserID:         user.ID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err = a.database.CreateOrganisationUser(c.Request().Context(), oru)
	if err != nil {
		return err
	}

	data := types.M{
		"status":  "success",
		"message": "Organisation created successfully",
		"data":    org,
	}

	return c.JSON(http.StatusCreated, data)
}
