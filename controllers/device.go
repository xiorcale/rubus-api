package controllers

import (
	"net/http"
	"os/exec"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/kjuvi/rubus-api/models"
	"github.com/kjuvi/rubus-api/services"
)

// Operations about Devices
type DeviceController struct {
	beego.Controller
}

// @Title AddDevice
// @Description Adds a `Device` into the database and prepare the necessary directory structure for deplo1ying it.
// @Param   hostname        query   string  true        "The hostname of the device"
// @Param	port			query	string	true		"The device's switch port"
// @Success 201 {object} models.Device
// @Failure 409 { "message": "conflict" }
// @Failure 500 { "message": "Internal Server Error" }
// @router /add [post]
func (d *DeviceController) AddDevice() {
	services.FilterAdmin(&d.Controller)

	hostname := d.GetString("hostname")
	port := d.GetString("port")

	// setup the necessary files and folders for the network boot and deployment
	cmd := exec.Command("./scripts/add-device.sh", hostname)
	go cmd.Run()

	// retrieve the device state
	device, jsonErr := services.GetDevice(port)
	if jsonErr != nil {
		d.Data["error"] = jsonErr
		d.Abort("JSONError")
	}

	// "cache" the device by inserting it into the database for faster read requests
	if jsonErr := models.AddDevice(device); jsonErr != nil {
		d.Data["error"] = jsonErr
		d.Abort("JSONError")
	}

	d.Ctx.Output.Status = http.StatusCreated
	d.Data["json"] = device
	d.ServeJSON()
}

// @Title DeleteDevice
// @Description Deletes a `Device` from the database and remove its directory structure used for deployment.
// @Param	deviceId		path 	int	true		"The device id to get"
// @Param   hostname        query   string  true        "The hostname of the device"
// @Success 204
// @Failure 409 { "message": "conflict" }
// @Failure 500 { "message": "Internal Server Error" }
// @router /:deviceId/delete [post]
func (d *DeviceController) DeleteDevice() {
	services.FilterAdmin(&d.Controller)
	hostname := d.GetString("hostname")
	deviceID, err := d.GetInt64(":deviceId")
	if err != nil {
		d.Data["error"] = models.NewBadRequestError
		d.Abort("JSONError")
	}

	// delete the necessary files and folders for the network boot and deployment
	cmd := exec.Command("./scripts/delete-device.sh", hostname)
	if err := cmd.Run(); err != nil {
		logs.Debug(err)
		d.Data["error"] = models.NewInternalServerError()
		d.Abort("JSONError")
	}

	// TODO: models.RemoveDevice() + return http status no content
	if jsonErr := models.DeleteDevice(deviceID); err != nil {
		d.Data["error"] = jsonErr
		d.Abort("JSONError")
	}

	d.Ctx.Output.Status = http.StatusOK
	d.ServeJSON()
}

// GetAll returns all the Rubus `Device`
// @Title GetAll
// @Description get all the rubus `Device`
// @Success 200 {object} []models.Device
// @Failure 500 { "message": "Internal Server Error" }
// @router / [get]
func (d *DeviceController) GetAll() {
	devices, jsonErr := models.GetAllDevices()
	if jsonErr != nil {
		d.Data["error"] = jsonErr
		d.Abort("JSONError")
	}

	d.Data["json"] = devices
	d.ServeJSON()
}

// Get a `Device`
// @Title Get
// @Description Get the Rubus `Device` with the given `deviceId`
// @Param	deviceId		path 	int	true		"The device id to get"
// @Success 200 {object} models.Device
// @Failure 400 { "message": "Bad Request Error" }
// @Failure 404 { "message": "User does not exists" }
// @Failure 500 { "message": "Internal Server Error" }
// @router /:deviceId [get]
func (d *DeviceController) Get() {
	deviceID, err := d.GetInt64(":deviceId")
	if err != nil {
		d.Data["error"] = models.NewBadRequestError
		d.Abort("JSONError")
	}

	device, jsonErr := models.GetDevice(deviceID)
	if jsonErr != nil {
		d.Data["error"] = jsonErr
		d.Abort("JSONError")
	}

	d.Ctx.Output.Status = http.StatusOK
	d.Data["json"] = device
	d.ServeJSON()
}

// PowerOn boots the `Device` on the given `port`
// @Title PowerOn
// @Description boots the `Device` on the given `port`. Note that on this context, the `deviceId` == device `port`.
// @Param	deviceId		path 	int	true		"The device port to turn on"
// @Success 204
// @router /:deviceId/on [post]
func (d *DeviceController) PowerOn() {
	port := d.GetString(":deviceId")
	deviceID, _ := strconv.Atoi(port)
	device, jsonErr := models.GetDevice(int64(deviceID))
	if jsonErr != nil {
		d.Data["error"] = jsonErr
		d.Abort("JSONError")
	}

	if device.Owner != nil {
		services.FilterOwnerOrAdmin(&d.Controller, *device.Owner)
	}

	if jsonErr := services.PowerDeviceOn(port); jsonErr != nil {
		d.Data["error"] = jsonErr
		d.Abort("JSONError")
	}

	models.SwitchDevicePower(device)

	d.Ctx.Output.Status = http.StatusNoContent
	d.ServeJSON()
}

// PowerOff shuts down the `Device` on the given `port`
// @Title PowerOff
// @Description shuts down the `Device` on the given `port`. Note that on this context, the `deviceId` == device `port`.
// @Param	deviceId		path 	int	true		"The device port to turn off"
// @Success 204
// @router /:deviceId/off [post]
func (d *DeviceController) PowerOff() {
	port := d.GetString(":deviceId")
	deviceID, _ := strconv.Atoi(port)
	device, jsonErr := models.GetDevice(int64(deviceID))
	if jsonErr != nil {
		d.Data["error"] = jsonErr
		d.Abort("JSONError")
	}

	if device.Owner != nil {
		services.FilterOwnerOrAdmin(&d.Controller, *device.Owner)
	}

	if jsonErr := services.PowerDeviceOff(port); jsonErr != nil {
		d.Data["error"] = jsonErr
		d.Abort("JSONError")
	}

	models.SwitchDevicePower(device)

	d.Ctx.Output.Status = http.StatusNoContent
	d.ServeJSON()
}
