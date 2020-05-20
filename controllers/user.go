package controllers

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo/v4"
	"github.com/xiorcale/rubus-api/models"
	"github.com/xiorcale/rubus-api/services"
	"gopkg.in/ini.v1"
)

// UserController -
type UserController struct {
	DB  *pg.DB
	Cfg *ini.File
}

// ListUser -
// @description Return a list containing all the `User`
// @id listUser
// @tags user
// @summary List all the users
// @produce json
// @security jwt
// @success 200 {array} models.User "A JSON array listing all the users"
// @router / [get]
func (u *UserController) ListUser(c echo.Context) error {
	if jsonErr := services.FilterAdmin(c); jsonErr != nil {
		return echo.NewHTTPError(jsonErr.Status, jsonErr)
	}

	users, jsonErr := models.GetAllUsers(u.DB)
	if jsonErr != nil {
		return echo.NewHTTPError(jsonErr.Status, jsonErr)
	}

	return c.JSON(http.StatusOK, users)
}

// GetMe -
// @description Return the `User` who made the request
// @id getMe
// @tags user
// @summary get the authenticated user
// @produce json
// @security jwt
// @success 200 {object} models.User "A JSON object describing a user"
// @router /me [get]
func (u *UserController) GetMe(c echo.Context) error {
	userID := ExtractIDFromToken(c)

	c.Logger().Printf("user id: ", userID)

	user, jsonErr := models.GetUser(u.DB, userID)
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
// @success 200 {object} models.User "AA JSON object describing a user"
// @router /me [put]
func (u *UserController) UpdateMe(c echo.Context) error {
	userID := ExtractIDFromToken(c)

	var user models.User
	cost, _ := u.Cfg.Section("security").Key("hashcost").Int()
	jsonErr := user.BindWithEmptyFields(c, cost)
	if jsonErr != nil {
		return echo.NewHTTPError(jsonErr.Status, jsonErr)
	}

	uu, jsonErr := models.UpdateUser(u.DB, userID, &user)
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
// @router /me [delete]
func (u *UserController) DeleteMe(c echo.Context) error {
	userID := ExtractIDFromToken(c)

	if jsonErr := models.DeleteUser(u.DB, userID); jsonErr != nil {
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
// @router /login [get]
func (u *UserController) Login(c echo.Context) error {
	username := c.QueryParam("username")
	password := c.QueryParam("password")

	uid, role, ok := models.Login(u.DB, username, password)
	if !ok {
		jsonErr := models.NewUnauthorizedError()
		return echo.NewHTTPError(jsonErr.Status, jsonErr)
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["ID"] = *uid
	claims["admin"] = (*role == models.EnumRoleAdmin)
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	secret := u.Cfg.Section("security").Key("jwtsecret").String()
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})

	// claims := &models.Claims{UserID: *uid, Role: *role}
	// token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)
	// secret := u.Cfg.Section("security").Key("jwtsecret").String()
	// tokenString, err := token.SignedString([]byte(secret))
	// if err != nil {
	// 	jsonErr := models.NewInternalServerError()
	// 	return echo.NewHTTPError(jsonErr.Status, jsonErr)
	// }

	// return c.JSON(http.StatusOK, map[string]string{"token": tokenString})
}
