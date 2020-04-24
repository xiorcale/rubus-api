package controllers

import (
	"net/http"
	"os/exec"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/kjuvi/rubus-api/models"
	"github.com/kjuvi/rubus-api/services"
)

// Operations about provisioning
type ProvisionerController struct {
	beego.Controller
}

// Acquire sets the request `User` as the owner of the `Device`
// @Title Acquire
// @Description sets the request `User` as the owner of the `Device`
// @Param	deviceId		path 	int	true		"The device port to acquire"
// @Success 200 {object} models.Device
// @router /:deviceId/acquire [post]
func (p *ProvisionerController) Acquire() {
	port, _ := p.GetInt64(":deviceId")
	claims := p.Ctx.Request.Context().Value("claims").(*models.Claims)

	// get the requested `Device`
	device, jsonErr := models.GetDevice(port)
	if jsonErr != nil {
		p.Data["error"] = jsonErr
		p.Abort("JSONError")
	}

	if device.Owner != nil {
		services.FilterOwnerOrAdmin(&p.Controller, *device.Owner)
	}

	if err := models.AcquireDevice(device, claims.UserID); err != nil {
		p.Data["error"] = models.NewInternalServerError()
		p.Abort("JSONError")
	}

	p.Ctx.Output.Status = http.StatusOK
	p.Data["json"] = device
	p.ServeJSON()
}

// Release removes the request `User` as the owner of the `Device`
// @Title Release
// @Description removes the request `User` as the owner of the `Device`
// @Param	deviceId		path 	int	true		"The device port to release"
// @Success 200 {object} models.Device
// @router /:deviceId/release [post]
func (p *ProvisionerController) Release() {
	port, _ := p.GetInt64(":deviceId")

	// get the requested `Device`
	device, jsonErr := models.GetDevice(port)
	if jsonErr != nil {
		p.Data["error"] = jsonErr
		p.Abort("JSONError")
	}

	if device.Owner != nil {
		services.FilterOwnerOrAdmin(&p.Controller, *device.Owner)
	}

	if err := models.ReleaseDevice(device); err != nil {
		p.Data["error"] = models.NewInternalServerError()
		p.Abort("JSONError")
	}

	p.Ctx.Output.Status = http.StatusOK
	p.Data["json"] = device
	p.ServeJSON()
}

// Deploy mounts the PXE boot folder into the tftp folder and reboot the `Device`
// @Title Deploy
// @Description mounts the PXE boot folder into the tftp folder and reboot the `Device`
// @Param	deviceId		path	int	true		"The device port to release"
// @Success	204
// @router /:deviceId/deploy [post]
func (p *ProvisionerController) Deploy() {
	port, _ := p.GetInt64(":deviceId")

	// get the requested `Device`
	device, jsonErr := models.GetDevice(port)
	if jsonErr != nil {
		p.Data["error"] = jsonErr
		p.Abort("JSONError")
	}

	if device.Owner != nil {
		services.FilterOwnerOrAdmin(&p.Controller, *device.Owner)
	}

	// setup the necessary files and folders for the network boot and deployment
	cmd := exec.Command("./scripts/deploy-device.sh", device.Hostname)
	go cmd.Run()

	if device.IsTurnOn {
		if jsonErr := services.PowerDeviceOff(strconv.FormatInt(port, 10)); jsonErr != nil {
			p.Data["error"] = jsonErr
			p.Abort("JSONError")
		}
		models.SwitchDevicePower(device)
	}

	if jsonErr := services.PowerDeviceOn(strconv.FormatInt(port, 10)); jsonErr != nil {
		p.Data["error"] = jsonErr
		p.Abort("JSONError")
	}

	models.SwitchDevicePower(device)

	p.Ctx.Output.Status = http.StatusNoContent
	p.ServeJSON()
}
