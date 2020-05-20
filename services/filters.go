package services

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/xiorcale/rubus-api/models"
)

// FilterAdmin checks if the `User` has an administrator `Role`. If not, return
// an Unauthorized `JSONError`.
func FilterAdmin(c echo.Context) *models.JSONError {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	isAdmin := claims["admin"].(bool)

	if !isAdmin {
		return models.NewUnauthorizedError()
	}

	return nil
}

// FilterMeOrAdmin checks if the `User` is the same as the given `uid` or and admin.
// If not, return an Unauthorized `JSONError`.
func FilterMeOrAdmin(c echo.Context, uid int64) *models.JSONError {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	isAdmin := claims["admin"].(bool)

	if claims["ID"] != uid && !isAdmin {
		return models.NewUnauthorizedError()
	}

	return nil
}

// FilterOwnerOrAdmin checks id the `User` is the same as the given `uid`
func FilterOwnerOrAdmin(c echo.Context, owner int64) *models.JSONError {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	isAdmin := claims["admin"].(bool)

	if claims["ID"] != owner && !isAdmin {
		return models.NewUnauthorizedError()
	}

	return nil
}
