package controllers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/xiorcale/rubus-api/models"
)

// ExtractIDFromToken extracts the `User` ID included inside the token from the
// current context
func ExtractIDFromToken(c echo.Context) int64 {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	id := claims["sub"].(float64)
	return int64(id)
}

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

// FilterIDOrAdmin checks if the `User` is the same as the given `uid` or and admin.
// If not, return an Unauthorized `JSONError`.
func FilterIDOrAdmin(c echo.Context, id int64) *models.JSONError {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	isAdmin := claims["admin"].(bool)
	requestID := ExtractIDFromToken(c)

	if requestID != id && !isAdmin {
		return models.NewUnauthorizedError()
	}

	return nil
}
