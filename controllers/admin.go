package controllers

import (
	"net/http"
	"os/exec"
	"strconv"

	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo/v4"
	"github.com/xiorcale/rubus-api/models"
	"github.com/xiorcale/rubus-api/services"
	"gopkg.in/ini.v1"
)

// AdminController -
type AdminController struct {
	DB  *pg.DB
	Cfg *ini.File
}

// CreateUser -
// @description Create a new Rubus `User` and save it into the database.
// @id createUser
// @tags admin
// @summary Create a new user
// @accept json
// @produce json
// @security jwt
// @param RequestBody body models.NewUser true "All the fields are required, except for the `role` which will default to `user` if not specified, and the expiration date which can be null."
// @success 201 {object} models.User
// @router /admin/user [post]
func (a *AdminController) CreateUser(c echo.Context) error {
	if jsonErr := FilterAdmin(c); jsonErr != nil {
		return echo.NewHTTPError(jsonErr.Status, jsonErr)
	}

	var user models.User
	cost, _ := a.Cfg.Section("security").Key("hashcost").Int()
	if jsonErr := user.Bind(c, cost); jsonErr != nil {
		return echo.NewHTTPError(jsonErr.Status, jsonErr)
	}

	if jsonErr := models.AddUser(a.DB, &user); jsonErr != nil {
		return echo.NewHTTPError(jsonErr.Status, jsonErr)
	}

	return c.JSON(http.StatusCreated, user)
}

// ListUser -
// @description Return a list containing all the `User`
// @id listUser
// @tags admin
// @summary List all the users
// @produce json
// @security jwt
// @success 200 {array} models.User "A JSON array listing all the users"
// @router /admin/user [get]
func (a *AdminController) ListUser(c echo.Context) error {
	if jsonErr := FilterAdmin(c); jsonErr != nil {
		return echo.NewHTTPError(jsonErr.Status, jsonErr)
	}

	users, jsonErr := models.GetAllUsers(a.DB)
	if jsonErr != nil {
		return echo.NewHTTPError(jsonErr.Status, jsonErr)
	}

	return c.JSON(http.StatusOK, users)
}

// CreateDevice -
// @description Add a `Device` into the database and prepare the necessary directory structure for deploying it.
// @id createDevice
// @tags admin
// @accept json
// @produce json
// @security jwt
// @param hostname query string true "The hostname of the device"
// @param port query string true "The device's switch port"
// @success 201 {object} models.Device
// @router /admin/device [post]
func (a *AdminController) CreateDevice(c echo.Context) error {
	if jsonErr := FilterAdmin(c); jsonErr != nil {
		return echo.NewHTTPError(jsonErr.Status, jsonErr)
	}

	hostname := c.QueryParam("hostname")
	port := c.QueryParam("port")

	// setup the necessary files and folders for the network boot and deployment
	cmd := exec.Command("./scripts/add-device.sh", hostname)
	go cmd.Run()

	// retrieve the device state
	device, jsonErr := services.GetDevice(port)
	if jsonErr != nil {
		return echo.NewHTTPError(jsonErr.Status, jsonErr)
	}

	// "cache" the device by inserting it into the
	// database for faster read requests
	if jsonErr := models.AddDevice(a.DB, device); jsonErr != nil {
		return echo.NewHTTPError(jsonErr.Status, jsonErr)
	}

	return c.JSON(http.StatusCreated, device)
}

// DeleteDevice -
// @description Delete a `Device` from the database and remove its directory structure used for deployment.
// @id deleteDevice
// @tags admin
// @summary Delete a device
// @produce json
// @security jwt
// @param hostname query string true "The hostname of the device"
// @param deviceId query int64 true "The device's switch port"
// @success 204
// @router /admin/device [delete]
func (a *AdminController) DeleteDevice(c echo.Context) error {
	if jsonErr := FilterAdmin(c); jsonErr != nil {
		return echo.NewHTTPError(jsonErr.Status, jsonErr)
	}

	hostname := c.QueryParam("hostname")
	deviceID, err := strconv.Atoi(c.QueryParam("deviceId"))
	if err != nil {
		jsonErr := models.NewBadRequestError()
		return echo.NewHTTPError(jsonErr.Status, jsonErr)
	}

	// delete the necessary files and folders
	// for the network boot and deployment
	cmd := exec.Command("./scripts/delete-device.sh", hostname)
	cmd.Run()

	if jsonErr := models.DeleteDevice(a.DB, int64(deviceID)); err != nil {
		return echo.NewHTTPError(jsonErr.Status, jsonErr)
	}

	return c.NoContent(http.StatusNoContent)
}
