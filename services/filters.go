package services

import (
	"github.com/xiorcale/rubus-api/models"
	"github.com/labstack/echo/v4"
)

// FilterAdmin checks if the `User` has an administrator `Role`. If not, return
// an Unauthorized `JSONError`.
func FilterAdmin(c echo.Context) *models.JSONError {
	claims := c.Get("claims").(*models.Claims)

	if claims.Role != models.EnumRoleAdmin {
		return models.NewUnauthorizedError()
	}

	return nil
}

// FilterMeOrAdmin checks if the `User` is the same as the given `uid` or and admin.
// If not, return an Unauthorized `JSONError`.
func FilterMeOrAdmin(c echo.Context, uid int64) *models.JSONError {
	claims := c.Get("claims").(*models.Claims)
	if claims.UserID != uid && claims.Role != models.EnumRoleAdmin {
		return models.NewUnauthorizedError()
	}

	return nil
}

// FilterOwnerOrAdmin checks id the `User` is the same as the given `uid`
func FilterOwnerOrAdmin(c echo.Context, owner int64) *models.JSONError {
	claims := c.Get("claims").(*models.Claims)
	if owner != claims.UserID && claims.Role != models.EnumRoleAdmin {
		return models.NewUnauthorizedError()
	}

	return nil
}
