package controllers

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/kjuvi/rubus-api/models"
	"github.com/kjuvi/rubus-api/services"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

// Operations about Users
type UserController struct {
	beego.Controller
}

// @Title ListUser
// @Description List all the `User`.
// @Success 200 {object} []models.User
// @Failure 500 { "message": "Internal Server Error" }
// @router / [get]
func (u *UserController) ListUser() {
	services.FilterAdmin(&u.Controller)
	users, jsonErr := models.GetAllUsers()
	if jsonErr != nil {
		u.Data["data"] = jsonErr
		u.Abort("JSONError")
	}

	u.Data["json"] = users
	u.ServeJSON()
}

// @Title GetMe
// @Description Return the `User` who made the request.
// @Success 200 {object} models.User
// @Failure 400 { "message": "Bad Request Error" }
// @Failure 404 { "message": "User does not exists" }
// @Failure 500 { "message": "Internal Server Error" }
// @router /me [get]
func (u *UserController) GetMe() {
	claims := u.Ctx.Request.Context().Value("claims").(*models.Claims)

	user, jsonErr := models.GetUser(claims.UserID)
	if jsonErr != nil {
		u.Data["error"] = jsonErr
		logs.Debug("JSON ERROR: ", jsonErr)
		u.Abort("JSONError")
	}

	u.Ctx.Output.Status = http.StatusOK
	u.Data["json"] = user
	u.ServeJSON()
}

// @Title UpdateMe
// @Description Update the `User` who made the request.
// @Param body body models.PutUser true "the `User` fields which can be updated. Giving all the fields is not mendatory, but at least one of them is required."
// @Success 200 {object} models.User
// @Failure 400 { "message": "Bad Request Error" }
// @Failure 404 { "message": "User does not exists" }
// @Failure 409 { "message": "conflict" }
// @Failure 500 { "message": "Internal Server Error" }
// @router /me [put]
func (u *UserController) UpdateMe() {
	claims := u.Ctx.Request.Context().Value("claims").(*models.Claims)

	var user models.User
	jsonErr := user.BindWithEmptyFields(u.Ctx.Input.RequestBody)
	if jsonErr != nil {
		u.Data["error"] = jsonErr
		u.Abort("JSONError")
	}

	uu, jsonErr := models.UpdateUser(claims.UserID, &user)
	if jsonErr != nil {
		u.Data["error"] = jsonErr
		u.Abort("JSONError")
	}

	u.Ctx.Output.Status = http.StatusOK
	u.Data["json"] = uu
	u.ServeJSON()
}

// @Title DeleteMe
// @Description Delete the `User` who made the request.
// @Success 200
// @Failure 400 { "message": "Bad Request Error" }
// @Failure 404 { "message": "User does not exists" }
// @Failure 500 { "message": "Internal Server Error" }
// @router /me [delete]
func (u *UserController) DeleteMe() {
	claims := u.Ctx.Request.Context().Value("claims").(*models.Claims)

	if jsonErr := models.DeleteUser(claims.UserID); jsonErr != nil {
		u.Data["error"] = jsonErr
		u.Abort("JSONError")
	}

	u.Ctx.Output.Status = http.StatusOK
	u.ServeJSON()
}

// @Title Login
// @Description Log a `User` into the system.
// @Param username query string true "The username used to login"
// @Param password query string true "The password used to login"
// @Success 200 { "token": "string" }
// @Failure 401 { "message": "Unauthorized" }
// @Failure 500 { "message": "Internal Server Error" }
// @router /login [get]
func (u *UserController) Login() {
	username := u.GetString("username")
	password := u.GetString("password")

	uid, role, ok := models.Login(username, password)
	if !ok {
		u.Data["error"] = models.NewUnauthorizedError()
		u.Abort("JSONError")
	}

	claims := &models.Claims{UserID: *uid, Role: *role}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)
	secret := beego.AppConfig.String("jwtsecret")
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		u.Data["error"] = models.NewInternalServerError
		u.Abort("JSONError")
	}

	u.Ctx.Output.Status = http.StatusOK
	u.Data["json"] = map[string]string{"token": tokenString}
	u.ServeJSON()
}
