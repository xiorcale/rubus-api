package controllers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

// ExtractIDFromToken extracts the `User` ID included inside the token from the
// current context
func ExtractIDFromToken(c echo.Context) int64 {
	token := c.Get("user").(*jwt.Token)
	c.Logger().Printf("token: ", token)
	claims := token.Claims.(jwt.MapClaims)
	id := claims["ID"].(float64)
	return int64(id)
}
