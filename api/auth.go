package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/oxiginedev/hng11-stage-two-task/pkg/datastore"
	"github.com/oxiginedev/hng11-stage-two-task/pkg/jwt"
	"github.com/oxiginedev/hng11-stage-two-task/pkg/models"
	"github.com/oxiginedev/hng11-stage-two-task/types"
)

type registerPayload struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Phone     string `json:"phone"`
}

// HandleRegister register a new user
func (a *API) HandleRegister(c echo.Context) error {
	var p registerPayload

	err := c.Bind(&p)
	if err != nil {
		return err
	}

	// validate user
	user := &models.User{
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
			return errors.New("account with phone or email exists")
		}
		return err
	}

	og := &models.Organisation{
		ID:          "",
		Name:        fmt.Sprintf("%s's Organisation", user.FirstName),
		Description: "Personal organisation",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = a.database.CreateOrganisation(c.Request().Context(), og)
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

type loginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// HandleLogin authenticates a user
func (a *API) HandleLogin(c echo.Context) error {
	var p loginPayload

	err := c.Bind(&p)
	if err != nil {
		return err
	}

	user, err := a.database.FetchUserByEmail(c.Request().Context(), p.Email)
	if err != nil {
		if errors.Is(err, datastore.ErrRecordNotFound) {
			return &APIError{StatusCode: http.StatusUnauthorized, Message: "Authentication failed"}
		}

		return err
	}

	ok, err := user.ComparePassword(p.Password)
	if err != nil {
		return err
	}

	if !ok {
		return &APIError{StatusCode: http.StatusUnauthorized, Message: "Authentication failed"}
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
