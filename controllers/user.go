package controllers

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/kjuvi/rubus-api/models"
	"github.com/kjuvi/rubus-api/services"

	"github.com/astaxie/beego"
)

// Operations about Users
type UserController struct {
	beego.Controller
}

// Post creates a new `User`
// @Title CreateUser
// @Description Creates a new Rubus `User` and save it into the database
// @Param	body		body 	models.NewUser	true		"body for user content"
// @Success 201 {object} models.User
// @Failure 409 { "message": "conflict" }
// @Failure 500 { "message": "Internal Server Error" }
// @router / [post]
func (u *UserController) Post() {
	services.FilterAdmin(&u.Controller)

	var user models.User
	if jsonErr := user.Bind(u.Ctx.Input.RequestBody); jsonErr != nil {
		u.Data["error"] = jsonErr
		u.Abort("JSONError")
	}

	if jsonErr := models.AddUser(&user); jsonErr != nil {
		u.Data["error"] = jsonErr
		u.Abort("JSONError")
	}

	u.Ctx.Output.Status = http.StatusCreated
	u.Data["json"] = user
	u.ServeJSON()
}

// GetAll returns all the Rubus `User`
// @Title GetAll
// @Description get all the rubus `User`
// @Success 200 {object} []models.User
// @Failure 500 { "message": "Internal Server Error" }
// @router / [get]
func (u *UserController) GetAll() {
	services.FilterAdmin(&u.Controller)
	users, jsonErr := models.GetAllUsers()
	if jsonErr != nil {
		u.Data["data"] = jsonErr
		u.Abort("JSONError")
	}

	u.Data["json"] = users
	u.ServeJSON()
}

// Get a `User`
// @Title Get
// @Description Get the authenticated Rubus `User`
// @Success 200 {object} models.User
// @Failure 400 { "message": "Bad Request Error" }
// @Failure 404 { "message": "User does not exists" }
// @Failure 500 { "message": "Internal Server Error" }
// @router /me [get]
func (u *UserController) Get() {
	claims := u.Ctx.Request.Context().Value("claims").(*models.Claims)

	user, jsonErr := models.GetUser(claims.UserID)
	if jsonErr != nil {
		u.Data["error"] = jsonErr
		u.Abort("JSONError")
	}

	u.Ctx.Output.Status = http.StatusOK
	u.Data["json"] = user
	u.ServeJSON()
}

// Put updates a user
// @Title Update
// @Description Update the Rubus `User` with the given `uid`
// @Param	body		body 	models.PutUser	true		"body for user content"
// @Success 200 {object} models.User
// @Failure 400 { "message": "Bad Request Error" }
// @Failure 404 { "message": "User does not exists" }
// @Failure 409 { "message": "conflict" }
// @Failure 500 { "message": "Internal Server Error" }
// @router /me [put]
func (u *UserController) Put() {
	claims := u.Ctx.Request.Context().Value("claims").(*models.Claims)

	var user models.User
	if jsonErr := user.BindWithEmptyFields(u.Ctx.Input.RequestBody); jsonErr != nil {
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

// Delete removes a `User`
// @Title Delete
// @Description delete the Rubus `User` with the given `uid`
// @Success 200
// @Failure 400 { "message": "Bad Request Error" }
// @Failure 404 { "message": "User does not exists" }
// @Failure 500 { "message": "Internal Server Error" }
// @router /me [delete]
func (u *UserController) Delete() {
	claims := u.Ctx.Request.Context().Value("claims").(*models.Claims)

	if jsonErr := models.DeleteUser(claims.UserID); jsonErr != nil {
		u.Data["error"] = jsonErr
		u.Abort("JSONError")
	}

	u.Ctx.Output.Status = http.StatusOK
	u.ServeJSON()
}

// Login logs a `User` into the system
// @Title Login
// @Description Logs user into the system
// @Param	username		query 	string	true		"The username for login"
// @Param	password		query 	string	true		"The password for login"
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
