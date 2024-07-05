package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/oxiginedev/hng11-stage-two-task/pkg/models"
)

var ErrTokenExpired = errors.New("token is expired")

type ValidatedToken struct {
	UserID string
	Expiry int64
}

func GenerateAccessToken(secret string, user *models.User, ttl int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Second * time.Duration(ttl)).Unix(),
	})

	accessToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func ValidateAccessToken(secret, accessToken string) (*ValidatedToken, error) {
	var userId string
	var expiry float64

	token, err := jwt.Parse(accessToken, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("jwt: unexpected signing method - %v", t.Header["alg"])
		}

		return []byte(secret), nil
	})

	if err != nil {
		v, ok := err.(*jwt.ValidationError)
		if ok && v.Errors == jwt.ValidationErrorExpired {
			if payload, ok := token.Claims.(jwt.MapClaims); ok {
				expiry = payload["exp"].(float64)
			}

			return &ValidatedToken{Expiry: int64(expiry)}, ErrTokenExpired
		}

		return nil, err
	}

	payload, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		userId = payload["sub"].(string)
		expiry = payload["exp"].(float64)

		vt := &ValidatedToken{UserID: userId, Expiry: int64(expiry)}

		return vt, nil
	}

	return nil, err
}
