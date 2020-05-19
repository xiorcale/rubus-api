package controllers

import (
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

// ExtractIDFromToken extracts the `User` ID included inside the token from the
// current context
func ExtractIDFromToken(c echo.Context) int64 {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	idStr := claims["ID"].(string)
	id, _ := strconv.Atoi(idStr)
	return int64(id)
}
