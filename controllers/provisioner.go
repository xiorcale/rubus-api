package controllers

import (
	"net/http"

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
