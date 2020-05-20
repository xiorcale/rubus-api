package controllers

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo/v4"
	"github.com/xiorcale/rubus-api/models"
	"gopkg.in/ini.v1"
)

// UserController -
type UserController struct {
	DB  *pg.DB
	Cfg *ini.File
}

// GetMe -
// @description Return the `User` who made the request
// @id getMe
// @tags user
// @summary get the authenticated user
// @produce json
// @security jwt
// @success 200 {object} models.User "A JSON object describing a user"
// @router /user/me [get]
func (u *UserController) GetMe(c echo.Context) error {
	id := ExtractIDFromToken(c)

	user, jsonErr := models.GetUser(u.DB, id)
	if jsonErr != nil {
		return echo.NewHTTPError(jsonErr.Status, jsonErr)
	}

	return c.JSON(http.StatusOK, user)
}

// UpdateMe -
// @description Update the `User` who made the request.
// @id updateMe
// @tags user
// @summary update the authenticated user
// @accept json
// @produce json
// @param RequestBody body models.PutUser true "the `User` fields which can be updated. Giving all the fields is not mendatory, but at least one of them is required."
// @success 200 {object} models.User "A JSON object describing a user"
// @router /user/me [put]
func (u *UserController) UpdateMe(c echo.Context) error {
	id := ExtractIDFromToken(c)

	var user models.User
	cost, _ := u.Cfg.Section("security").Key("hashcost").Int()
	jsonErr := user.BindWithEmptyFields(c, cost)
	if jsonErr != nil {
		return echo.NewHTTPError(jsonErr.Status, jsonErr)
	}

	uu, jsonErr := models.UpdateUser(u.DB, id, &user)
	if jsonErr != nil {
		return echo.NewHTTPError(jsonErr.Status, jsonErr)
	}

	return c.JSON(http.StatusOK, uu)
}

// DeleteMe -
// @description Delete the `User` who made the request.
// @id deleteMe
// @tags user
// @summary delethe the autenticated user
// @produce json
// @success 200
// @router /user/me [delete]
func (u *UserController) DeleteMe(c echo.Context) error {
	id := ExtractIDFromToken(c)

	if jsonErr := models.DeleteUser(u.DB, id); jsonErr != nil {
		return echo.NewHTTPError(jsonErr.Status, jsonErr)
	}

	return c.NoContent(http.StatusOK)
}

// Login -
// @description Log a `User` into the system.
// @id login
// @tags user
// @summary Log a user in
// @accept json
// @produce json
// @param username query string true "The username used to login"
// @param password query string true "The password used to login"
// @success 200
// @router /user/login [get]
func (u *UserController) Login(c echo.Context) error {
	username := c.QueryParam("username")
	password := c.QueryParam("password")

	user := models.Login(u.DB, username, password)

	if user == nil || (user.Expiration.Unix() > 0 && user.Expiration.Before(time.Now())) {
		jsonErr := models.NewUnauthorizedError()
		return echo.NewHTTPError(jsonErr.Status, jsonErr)
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["ID"] = user.ID
	claims["admin"] = (user.Role == models.EnumRoleAdmin)

	if user.Expiration.Unix() > 0 {
		claims["exp"] = user.Expiration.Unix()
	}

	secret := u.Cfg.Section("security").Key("jwtsecret").String()
	t, _ := token.SignedString([]byte(secret))

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}
