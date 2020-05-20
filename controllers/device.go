package controllers

import (
	"net/http"
	"strconv"

	"github.com/go-pg/pg/v9"
	"github.com/xiorcale/rubus-api/models"
	"github.com/xiorcale/rubus-api/services"
	"github.com/labstack/echo/v4"
)

// DeviceController -
type DeviceController struct {
	DB *pg.DB
}

// ListDevice -
// @description List all the `Device`
// @id listDevice
// @tags device
// @summary list all the devices
// @produce json
// @security jwt
// @success 200 {array} models.Device "A JSON array listing all the devices"
// @router / [get]
func (d *DeviceController) ListDevice(c echo.Context) error {
	devices, jsonErr := models.GetAllDevices(d.DB)
	if jsonErr != nil {
		return echo.NewHTTPError(jsonErr.Status, jsonErr)
	}

	return c.JSON(http.StatusOK, devices)
}

// Get -
// @description Return the `Device` with the given `deviceId`.
// @id getDevice
// @tags device
// @summary get a device by id
// @produce json
// @security jwt
// @param deviceId path int true "The id of the `Device` to get"
// @success 200 {object} models.Device
// @router /:deviceId [get]
func (d *DeviceController) Get(c echo.Context) error {
	deviceID, err := strconv.Atoi(c.Param("deviceId"))

	if err != nil {
		jsonErr := models.NewBadRequestError()
		return echo.NewHTTPError(jsonErr.Status, jsonErr)
	}

	device, jsonErr := models.GetDevice(d.DB, int64(deviceID))
	if jsonErr != nil {
		return echo.NewHTTPError(jsonErr.Status, jsonErr)
	}

	return c.JSON(http.StatusOK, device)
}

// PowerOn -
// @description Boot the `Device` with the given `deviceId`.
// @id powerOn
// @tags device
// @summary Boot a device
// @produce json
// @security jwt
// @param deviceId path int true "The device id to turn on"
// @success 204
// @router /:deviceId/on [post]
func (d *DeviceController) PowerOn(c echo.Context) error {
	port := c.Param("deviceId")
	deviceID, _ := strconv.Atoi(port)
	device, jsonErr := models.GetDevice(d.DB, int64(deviceID))
	if jsonErr != nil {
		return echo.NewHTTPError(jsonErr.Status, jsonErr)
	}

	if device.Owner != nil {
		if jsonErr := services.FilterOwnerOrAdmin(c, *device.Owner); jsonErr != nil {
			return echo.NewHTTPError(jsonErr.Status, jsonErr)
		}
	}

	if jsonErr := services.PowerDeviceOn(port); jsonErr != nil {
		return echo.NewHTTPError(jsonErr.Status, jsonErr)
	}

	if !device.IsTurnOn {
		models.SwitchDevicePower(d.DB, device)
	}

	return c.NoContent(http.StatusNoContent)
}

// PowerOff -
// @description Shuts down the `Device` on the given `port`
// @id powerOff
// @tags device
// @summary Shut down a device
// @produce json
// @security jwt
// @param deviceId path int true "The device id to turn off"
// @success 204
// @router /:deviceId/off [post]
func (d *DeviceController) PowerOff(c echo.Context) error {
	port := c.Param("deviceId")
	deviceID, _ := strconv.Atoi(port)
	device, jsonErr := models.GetDevice(d.DB, int64(deviceID))
	if jsonErr != nil {
		return echo.NewHTTPError(jsonErr.Status, jsonErr)
	}

	if device.Owner != nil {
		if jsonErr := services.FilterOwnerOrAdmin(c, *device.Owner); jsonErr != nil {
			return echo.NewHTTPError(jsonErr.Status, jsonErr)
		}
	}

	if jsonErr := services.PowerDeviceOff(port); jsonErr != nil {
		return echo.NewHTTPError(jsonErr.Status, jsonErr)
	}

	if device.IsTurnOn {
		models.SwitchDevicePower(d.DB, device)
	}

	return c.NoContent(http.StatusNoContent)
}
