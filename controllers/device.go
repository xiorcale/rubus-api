package controllers

import (
	"net/http"
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

// RegisterAll adds all `Device` from the provider into the database
// @Title RegisterAll
// @Description Registers all the `Device` from the provider into the database
// @Success 201 {object} []models.Device
// @Failure 409 { "message": "conflict" }
// @Failure 500 { "message": "Internal Server Error" }
// @router /register [post]
func (d *DeviceController) RegisterAll() {
	devices, err := services.GetAllDevices()
	if err != nil {
		logs.Debug(err)
		d.Data["status"] = http.StatusInternalServerError
		d.Data["msg"] = err.Error()
		d.Abort("JSONError")
	}

	if err := models.AddDeviceMulti(devices); err != nil {
		if err.Error() == "ID already exists" {
			d.Data["status"] = http.StatusConflict
		} else {
			d.Data["status"] = http.StatusInternalServerError
		}
		d.Data["msg"] = err.Error()
		d.Abort("JSONError")
	}

	d.Ctx.Output.Status = http.StatusCreated
	d.Data["json"] = devices
	d.ServeJSON()
}

// Register adds a `Device` from the provider into the database
// @Title Register
// @Description Registers a `Device` from the provider into the database. Note that on this context, the `deviceId` == device `port`.
// @Param	deviceId		path 	int	true		"The device id to register"
// @Success 201 {object} models.Device
// @Failure 409 { "message": "conflict" }
// @Failure 500 { "message": "Internal Server Error" }
// @router /:deviceId/register [post]
func (d *DeviceController) Register() {
	port := d.GetString(":deviceId")

	device, err := services.GetDevice(port)
	if err != nil {
		logs.Debug(err)
		d.Data["status"] = http.StatusInternalServerError
		d.Data["msg"] = err.Error()
		d.Abort("JSONError")
	}

	if err := models.AddDevice(device); err != nil {
		if err.Error() == "ID already exists" {
			d.Data["status"] = http.StatusConflict
		} else {
			d.Data["status"] = http.StatusInternalServerError
		}
		d.Data["msg"] = err.Error()
		d.Abort("JSONError")
	}

	d.Ctx.Output.Status = http.StatusCreated
	d.Data["json"] = device
	d.ServeJSON()
}

// GetAll returns all the Rubus `Device`
// @Title GetAll
// @Description get all the rubus `Device`
// @Success 200 {object} []models.Device
// @Failure 500 { "message": "Internal Server Error" }
// @router / [get]
func (d *DeviceController) GetAll() {
	devices, err := models.GetAllDevices()
	if err != nil {
		d.Data["status"] = http.StatusInternalServerError
		d.Data["msg"] = err.Error()
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
		d.Data["status"] = http.StatusBadRequest
		d.Data["msg"] = "Bad Request Error"
		d.Abort("JSONError")
	}

	device, err := models.GetDevice(deviceID)
	if err != nil {
		if err.Error() == "Device does not exists" {
			d.Data["status"] = http.StatusNotFound
		} else {
			d.Data["status"] = http.StatusInternalServerError
		}
		d.Data["msg"] = err.Error()
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

	if err := services.PowerDeviceOn(port); err != nil {
		logs.Debug(err)
		d.Data["status"] = http.StatusInternalServerError
		d.Data["msg"] = err.Error()
		d.Abort("JSONError")
	}

	deviceID, _ := strconv.Atoi(port)
	models.SwitchDevicePower(int64(deviceID))

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

	if err := services.PowerDeviceOff(port); err != nil {
		logs.Debug(err)
		d.Data["status"] = http.StatusInternalServerError
		d.Data["msg"] = err.Error()
		d.Abort("JSONError")
	}

	deviceID, _ := strconv.Atoi(port)
	models.SwitchDevicePower(int64(deviceID))

	d.Ctx.Output.Status = http.StatusNoContent
	d.ServeJSON()
}
