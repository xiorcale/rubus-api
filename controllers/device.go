package controllers

import (
	"net/http"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/kjuvi/rubus-api/models"
	"github.com/kjuvi/rubus-api/services"
)

// Operations about devices such as provisioning or deployment
type DeviceController struct {
	beego.Controller
}

// @Title ListDevice
// @Description List all the `Device`.
// @Success 200 {object} []models.Device
// @Failure 500 { "message": "Internal Server Error" }
// @router / [get]
func (d *DeviceController) ListDevice() {
	devices, jsonErr := models.GetAllDevices()
	if jsonErr != nil {
		d.Data["error"] = jsonErr
		d.Abort("JSONError")
	}

	d.Data["json"] = devices
	d.ServeJSON()
}

// @Title GetDevice
// @Description Return the `Device` with the given `deviceId`.
// @Param deviceId path int true "The id of the `Device` to get"
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

// @Title PowerOn
// @Description Boot the `Device` with the given `deviceId`.
// @Param deviceId path int true "The device id to turn on"
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

	if !device.IsTurnOn {
		models.SwitchDevicePower(device)
	}

	d.Ctx.Output.Status = http.StatusNoContent
	d.ServeJSON()
}

// PowerOff shuts down the `Device` on the given `port`
// @Title PowerOff
// @Description Shut down the `Device` with the given `deviceId`.
// @Param deviceId path int true "The device id to turn off"
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

	if device.IsTurnOn {
		models.SwitchDevicePower(device)
	}

	d.Ctx.Output.Status = http.StatusNoContent
	d.ServeJSON()
}
