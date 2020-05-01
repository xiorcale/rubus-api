package controllers

import (
	"net/http"
	"os/exec"

	"github.com/astaxie/beego"
	"github.com/kjuvi/rubus-api/models"
	"github.com/kjuvi/rubus-api/services"
)

// Operations which require administrative rights
type AdminController struct {
	beego.Controller
}

// @Title CreateUser
// @Description Create a new Rubus `User` and save it into the database.
// @Param body body models.NewUser true "All the fields are required, except for the `role` which will default to `user` if not specified."
// @Success 201 {object} models.User
// @Failure 409 { "message": "conflict" }
// @Failure 500 { "message": "Internal Server Error" }
// @router /user [post]
func (a *AdminController) CreateUser() {
	services.FilterAdmin(&a.Controller)

	var user models.User
	if jsonErr := user.Bind(a.Ctx.Input.RequestBody); jsonErr != nil {
		a.Data["error"] = jsonErr
		a.Abort("JSONError")
	}

	if jsonErr := models.AddUser(&user); jsonErr != nil {
		a.Data["error"] = jsonErr
		a.Abort("JSONError")
	}

	a.Ctx.Output.Status = http.StatusCreated
	a.Data["json"] = user
	a.ServeJSON()
}

// @Title AddDevice
// @Description Add a `Device` into the database and prepare the necessary directory structure for deploying it.
// @Param hostname query string true "The hostname of the device"
// @Param port query string true "The device's switch port"
// @Success 201 {object} models.Device
// @Failure 409 { "message": "conflict" }
// @Failure 500 { "message": "Internal Server Error" }
// @router /device [post]
func (a *AdminController) CreateDevice() {
	services.FilterAdmin(&a.Controller)

	hostname := a.GetString("hostname")
	port := a.GetString("port")

	// setup the necessary files and folders for the network boot and deployment
	cmd := exec.Command("./scripts/add-device.sh", hostname)
	go cmd.Run()

	// retrieve the device state
	device, jsonErr := services.GetDevice(port)
	if jsonErr != nil {
		a.Data["error"] = jsonErr
		a.Abort("JSONError")
	}

	// "cache" the device by inserting it into the
	// database for faster read requests
	if jsonErr := models.AddDevice(device); jsonErr != nil {
		a.Data["error"] = jsonErr
		a.Abort("JSONError")
	}

	a.Ctx.Output.Status = http.StatusCreated
	a.Data["json"] = device
	a.ServeJSON()
}

// @Title DeleteDevice
// @Description Delete a `Device` from the database and remove its directory structure used for deployment.
// @Param hostname query string true "The hostname of the device"
// @Param deviceId query int64 true "The device's switch port"
// @Success 204
// @Failure 400 { "message": "Bad Request Error" }
// @Failure 404 { "message": "Not Found" }
// @Failure 500 { "message": "Internal Server Error" }
// @router /device [delete]
func (a *AdminController) DeleteDevice() {
	services.FilterAdmin(&a.Controller)

	hostname := a.GetString("hostname")
	deviceID, err := a.GetInt64("deviceId")
	if err != nil {
		a.Data["error"] = models.NewBadRequestError
		a.Abort("JSONError")
	}

	// delete the necessary files and folders
	// for the network boot and deployment
	cmd := exec.Command("./scripts/delete-device.sh", hostname)
	cmd.Run()

	if jsonErr := models.DeleteDevice(deviceID); err != nil {
		a.Data["error"] = jsonErr
		a.Abort("JSONError")
	}

	a.Ctx.Output.Status = http.StatusNoContent
	a.ServeJSON()
}
