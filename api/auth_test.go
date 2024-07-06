package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandleRegister(t *testing.T) {
	// Create a new instance of the API
	api := &API{}

	// Create a new Echo instance
	e := echo.New()

	// Create a RegisterPayload with test data
	payload := RegisterPayload{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Password:  "password",
		Phone:     "1234567890",
	}

	// Convert the payload to JSON
	payloadJSON, _ := json.Marshal(payload)

	// Create a new HTTP request with the JSON payload
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(payloadJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Create a new HTTP response recorder
	rec := httptest.NewRecorder()

	// Set up the Echo context
	c := e.NewContext(req, rec)

	// Call the HandleRegister function
	err := api.HandleRegister(c)

	// Assert that no error occurred
	assert.NoError(t, err)

	// Assert that the response status code is 201 (Created)
	assert.Equal(t, http.StatusCreated, rec.Code)

	// Assert that the response body contains the expected message
	expectedMessage := "User registered successfully"
	assert.Equal(t, expectedMessage, rec.Body.String())
}
