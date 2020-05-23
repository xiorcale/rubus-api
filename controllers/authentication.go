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

// AuthenticationController -
type AuthenticationController struct {
    DB *pg.DB
    Cfg *ini.File
}


// Login -
// @description Log a `User` into the system.
// @id login
// @tags authentication
// @summary Log a user in
// @accept json
// @produce json
// @param username query string true "The username used to login"
// @param password query string true "The password used to login"
// @success 200
// @router /login [get]
func (a *AuthenticationController) Login(c echo.Context) error {
    username := c.QueryParam("username")
    password := c.QueryParam("password")

    user := models.Login(a.DB, username, password)

    if user == nil || (user.Expiration.Unix() > 0 && user.Expiration.Before(time.Now())) {
        jsonErr := models.NewUnauthorizedError()
        return echo.NewHTTPError(jsonErr.Status, jsonErr)
    }

    token := jwt.New(jwt.SigningMethodHS256)

    claims := token.Claims.(jwt.MapClaims)
    claims["sub"] = user.ID
    claims["admin"] = (user.Role == models.EnumRoleAdmin)

    if user.Expiration.Unix() > 0 {
        claims["exp"] = user.Expiration.Unix()
    }

    secret := a.Cfg.Section("security").Key("jwtsecret").String()
    t, _ := token.SignedString([]byte(secret))

    return c.JSON(http.StatusOK, map[string]string{
        "token": t,
    })
}
