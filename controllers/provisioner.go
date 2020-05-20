package controllers

import (
	"net/http"
	"os/exec"
	"strconv"

	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo/v4"
	"github.com/xiorcale/rubus-api/models"
	"github.com/xiorcale/rubus-api/services"
)

// ProvisionerController -
type ProvisionerController struct {
	DB *pg.DB
}

// Acquire -
// @description Set the `User` who made the request as the owner of the `Device`.
// @id acquire
// @tags device
// @summary acquire a device
// @produce json
// @security jwt
// @param deviceId path int true "The id of the `Device` to acquire"
// @success 200 {object} models.Device
// @router /:deviceId/acquire [post]
func (p *ProvisionerController) Acquire(c echo.Context) error {
	port, _ := strconv.Atoi(c.Param("deviceId"))
	userID := ExtractIDFromToken(c)

	// get the requested `Device`
	device, jsonErr := models.GetDevice(p.DB, int64(port))
	if jsonErr != nil {
		return echo.NewHTTPError(jsonErr.Status, jsonErr)
	}

	if device.Owner != nil {
		services.FilterOwnerOrAdmin(c, *device.Owner)
	}

	if err := models.AcquireDevice(p.DB, device, userID); err != nil {
		return echo.NewHTTPError(jsonErr.Status, jsonErr)
	}

	return c.JSON(http.StatusOK, device)
}

// Release -
// @description Remove the `Device`'s ownership from the `User` who made the request.
// @id release
// @tags device
// @summary release a device
// @produce json
// @security jwt
// @param	deviceId		path 	int	true		"The device port to release"
// @success 200 {object} models.Device
// @router /:deviceId/release [post]
func (p *ProvisionerController) Release(c echo.Context) error {
	port, _ := strconv.Atoi(c.Param("deviceId"))

	// get the requested `Device`
	device, jsonErr := models.GetDevice(p.DB, int64(port))
	if jsonErr != nil {
		return echo.NewHTTPError(jsonErr.Status, jsonErr)
	}

	if device.Owner != nil {
		if jsonErr := services.FilterOwnerOrAdmin(c, *device.Owner); jsonErr != nil {
			return echo.NewHTTPError(jsonErr.Status, jsonErr)
		}
	}

	if err := models.ReleaseDevice(p.DB, device); err != nil {
		return echo.NewHTTPError(jsonErr.Status, jsonErr)
	}

	return c.JSON(http.StatusOK, device)
}

// Deploy -
// @description Configure the PXE boot for the `Device` and reboot it.
// @id deploy
// @tags device
// @summary deploy a device
// @produce json
// @security jwt
// @param deviceId path int true "The device id to deploy"
// @success	204
// @router /:deviceId/deploy [post]
func (p *ProvisionerController) Deploy(c echo.Context) error {
	port, _ := strconv.Atoi(c.Param("deviceId"))

	// get the requested `Device`
	device, jsonErr := models.GetDevice(p.DB, int64(port))
	if jsonErr != nil {
		return echo.NewHTTPError(jsonErr.Status, jsonErr)
	}

	if device.Owner != nil {
		if jsonErr := services.FilterOwnerOrAdmin(c, *device.Owner); jsonErr != nil {
			return echo.NewHTTPError(jsonErr.Status, jsonErr)
		}
	}

	// setup the necessary files and folders for the network boot and deployment
	cmd := exec.Command("./scripts/deploy-device.sh", device.Hostname)
	go cmd.Run()

	if device.IsTurnOn {
		jsonErr := services.PowerDeviceOff(strconv.FormatInt(int64(port), 10))
		if jsonErr != nil {
			return echo.NewHTTPError(jsonErr.Status, jsonErr)
		}
		models.SwitchDevicePower(p.DB, device)
	}

	jsonErr = services.PowerDeviceOn(strconv.FormatInt(int64(port), 10))
	if jsonErr != nil {
		return echo.NewHTTPError(jsonErr.Status, jsonErr)
	}

	models.SwitchDevicePower(p.DB, device)

	return c.NoContent(http.StatusNoContent)
}
