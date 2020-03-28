package controllers

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/kjuvi/rubus-api/models"

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
// @Description Get the Rubus `User` with the given `uid`
// @Param	uid		path 	int	true		"The user id to get"
// @Success 200 {object} models.User
// @Failure 400 { "message": "Bad Request Error" }
// @Failure 404 { "message": "User does not exists" }
// @Failure 500 { "message": "Internal Server Error" }
// @router /:uid [get]
func (u *UserController) Get() {
	uid, err := u.GetInt64(":uid")
	if err != nil {
		u.Data["error"] = models.NewBadRequestError()
		u.Abort("JSONError")
	}

	user, jsonErr := models.GetUser(int64(uid))
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
// @Param	uid		path 	int 	true		"The user id to update"
// @Param	body		body 	models.NewUser	true		"body for user content"
// @Success 200 {object} models.User
// @Failure 400 { "message": "Bad Request Error" }
// @Failure 404 { "message": "User does not exists" }
// @Failure 409 { "message": "conflict" }
// @Failure 500 { "message": "Internal Server Error" }
// @router /:uid [put]
func (u *UserController) Put() {
	uid, err := u.GetInt64(":uid")
	if err != nil {
		u.Data["error"] = models.NewBadRequestError()
		u.Abort("JSONError")
	}

	var user models.User
	if jsonErr := user.Bind(u.Ctx.Input.RequestBody); jsonErr != nil {
		u.Data["error"] = jsonErr
		u.Abort("JSONError")
	}

	uu, jsonErr := models.UpdateUser(uid, &user)
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
// @Param	uid		path 	string	true		"The user id to delete"
// @Success 200
// @Failure 400 { "message": "Bad Request Error" }
// @Failure 404 { "message": "User does not exists" }
// @Failure 500 { "message": "Internal Server Error" }
// @router /:uid [delete]
func (u *UserController) Delete() {
	uid, err := u.GetInt64(":uid")
	if err != nil {
		u.Data["status"] = http.StatusBadRequest
		u.Data["msg"] = "Bad Request Error"
		u.Abort("JSONError")
	}

	if jsonErr := models.DeleteUser(uid); jsonErr != nil {
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

	uid, ok := models.Login(username, password)
	if !ok {
		u.Data["error"] = models.NewUnauthorizedError()
		u.Abort("JSONError")
	}

	tk := &models.Token{UserID: *uid}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
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
