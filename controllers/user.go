package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

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
// @router / [post]
func (u *UserController) Post() {
	var user models.User
	if err := user.Bind(u.Ctx.Input.RequestBody); err != nil {
		u.Data["status"] = http.StatusBadRequest
		u.Data["msg"] = err.Error()
		u.Abort("JSONError")
	}

	if err := models.AddUser(&user); err != nil {
		if strings.Contains(err.Error(), "username") {
			u.Data["status"] = http.StatusConflict
		} else {
			u.Data["status"] = http.StatusInternalServerError
		}
		u.Data["msg"] = err.Error()
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
// @router / [get]
func (u *UserController) GetAll() {
	users, err := models.GetAllUsers()
	if err != nil {
		u.Data["status"] = http.StatusInternalServerError
		u.Data["msg"] = err.Error()
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
// @router /:uid [get]
func (u *UserController) Get() {
	uid, err := u.GetInt64(":uid")
	if err != nil {
		u.Data["status"] = http.StatusBadRequest
		u.Data["msg"] = "Bad Request Error"
		u.Abort("JSONError")
	}

	user, err := models.GetUser(uid)
	if err != nil {
		if err.Error() == "User does not exists" {
			u.Data["status"] = http.StatusNotFound
		} else {
			u.Data["status"] = http.StatusInternalServerError
		}
		u.Data["msg"] = err.Error()
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
// @Failure 409 { "message": "conflict" }
// @router /:uid [put]
func (u *UserController) Put() {
	uid, err := u.GetInt64(":uid")
	if err != nil {
		u.Data["status"] = http.StatusBadRequest
		u.Data["msg"] = "Bad Request Error"
		u.Abort("JSONError")
	}

	var user models.User
	if err := user.Bind(u.Ctx.Input.RequestBody); err != nil {
		u.Data["status"] = http.StatusBadRequest
		u.Data["msg"] = err.Error()
		u.Abort("JSONError")
	}

	uu, err := models.UpdateUser(uid, &user)
	if err != nil {
		if strings.Contains(err.Error(), "username") {
			u.Data["status"] = http.StatusConflict
		} else if err.Error() == "User does not exists" {
			u.Data["status"] = http.StatusNotFound
		} else {
			u.Data["status"] = http.StatusInternalServerError
		}
		u.Data["msg"] = err.Error()
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
// @router /:uid [delete]
func (u *UserController) Delete() {
	uid, err := u.GetInt64(":uid")
	if err != nil {
		u.Data["status"] = http.StatusBadRequest
		u.Data["msg"] = "Bad Request Error"
		u.Abort("JSONError")
	}

	if err := models.DeleteUser(uid); err != nil {
		if err.Error() == "User does not exists" {
			u.Data["status"] = http.StatusNotFound
		} else {
			u.Data["status"] = http.StatusInternalServerError
		}
		u.Data["msg"] = err.Error()
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
// @router /login [get]
func (u *UserController) Login() {
	username := u.GetString("username")
	password := u.GetString("password")

	uid, ok := models.Login(username, password)
	if !ok {
		u.Data["status"] = http.StatusUnauthorized
		u.Data["msg"] = "Unauthorized"
		u.Abort("JSONError")
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = strconv.FormatInt(*uid, 10)
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	secret := beego.AppConfig.String("jwtsecret")
	t, err := token.SignedString([]byte(secret)) // TODO: don't hard coded JWT_SECRET
	if err != nil {
		u.Data["status"] = http.StatusInternalServerError
		u.Data["msg"] = "Internal Server Error"
		u.Abort("JSONError")
	}

	u.Ctx.Output.Status = http.StatusOK
	u.Data["json"] = map[string]string{"token": t}
	u.ServeJSON()
}
