package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/oklog/ulid/v2"
	"github.com/oxiginedev/hng11-stage-two-task/pkg/datastore"
	"github.com/oxiginedev/hng11-stage-two-task/pkg/jwt"
	"github.com/oxiginedev/hng11-stage-two-task/pkg/models"
	"github.com/oxiginedev/hng11-stage-two-task/types"
)

type RegisterPayload struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
	Phone     string `json:"phone" validate:"required,numeric"`
}

// HandleRegister register a new user
func (a *API) HandleRegister(c echo.Context) error {
	var p RegisterPayload

	err := c.Bind(&p)
	if err != nil {
		return err
	}

	if err := c.Validate(p); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			log.Print(ve)
			return &APIError{
				Errors: ve,
			}
		}
		return err
	}

	user := &models.User{
		ID:        ulid.Make().String(),
		FirstName: p.FirstName,
		LastName:  p.LastName,
		Email:     p.Email,
		Phone:     p.Phone,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = user.SetPassword(p.Password)
	if err != nil {
		return err
	}

	err = a.database.CreateUser(c.Request().Context(), user)
	if err != nil {
		if errors.Is(err, datastore.ErrDuplicate) {
			return newAPIError(http.StatusUnprocessableEntity, "Account with email or phone exists")
		}
		return err
	}

	og := &models.Organisation{
		ID:          ulid.Make().String(),
		Name:        fmt.Sprintf("%s's Organisation", user.FirstName),
		Description: "This is an organisation",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = a.database.CreateOrganisation(c.Request().Context(), og)
	if err != nil {
		return err
	}

	ogu := &models.OrganisationUser{
		ID:             ulid.Make().String(),
		OrganisationID: og.ID,
		UserID:         user.ID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err = a.database.CreateOrganisationUser(c.Request().Context(), ogu)
	if err != nil {
		return err
	}

	accessToken, err := jwt.GenerateAccessToken(a.config.JWT.Secret, user, a.config.JWT.Expiry)
	if err != nil {
		return err
	}

	data := types.M{
		"status":  "success",
		"message": "Registration successful",
		"data": types.M{
			"accessToken": accessToken,
			"user":        user,
		},
	}

	return c.JSON(http.StatusCreated, data)
}

type LoginPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// HandleLogin authenticates a user
func (a *API) HandleLogin(c echo.Context) error {
	var p LoginPayload

	err := c.Bind(&p)
	if err != nil {
		return err
	}

	if err := c.Validate(p); err != nil {
		return err
	}

	user, err := a.database.FetchUserByEmail(c.Request().Context(), p.Email)
	if err != nil {
		if errors.Is(err, datastore.ErrRecordNotFound) {
			return newAPIError(http.StatusUnauthorized, "Authentication failed")
		}

		return err
	}

	ok, err := user.ComparePassword(p.Password)
	if err != nil {
		return err
	}

	if !ok {
		return newAPIError(http.StatusUnauthorized, "Authentication failed")
	}

	accessToken, err := jwt.GenerateAccessToken(a.config.JWT.Secret, user, 1800)
	if err != nil {
		return err
	}

	data := types.M{
		"status":  "success",
		"message": "Login successful",
		"data": types.M{
			"accessToken": accessToken,
			"user":        user,
		},
	}

	return c.JSON(http.StatusOK, data)
}
