package controllers

import (
	"net/http"

	"github.com/astaxie/beego"
	"github.com/kjuvi/rubus-api/models"
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
	uid := p.Ctx.Request.Context().Value("user").(int64)

	// get the requested `Device`
	device, err := models.GetDevice(port)
	if err != nil {
		if err.Error() == "Device does not exists" {
			p.Data["status"] = http.StatusNotFound
		} else {
			p.Data["status"] = http.StatusInternalServerError
		}
		p.Data["msg"] = err.Error()
		p.Abort("JSONError")
	}

	// check the `Device` is unowned
	if device.Owner != nil {
		p.Data["status"] = http.StatusForbidden
		p.Data["msg"] = "Forbidden"
		p.Abort("JSONError")
	}

	if err := models.AcquireDevice(device, uid); err != nil {
		p.Data["status"] = http.StatusInternalServerError
		p.Data["msg"] = err.Error()
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
	uid := p.Ctx.Request.Context().Value("user").(int64)

	// get the requested `Device`
	device, err := models.GetDevice(port)
	if err != nil {
		if err.Error() == "Device does not exists" {
			p.Data["status"] = http.StatusNotFound
		} else {
			p.Data["status"] = http.StatusInternalServerError
		}
		p.Data["msg"] = err.Error()
		p.Abort("JSONError")
	}

	if *device.Owner != uid {
		p.Data["status"] = http.StatusForbidden
		p.Data["msg"] = "Forbidden"
		p.Abort("JSONError")
	}

	if err := models.ReleaseDevice(device); err != nil {
		p.Data["status"] = http.StatusInternalServerError
		p.Data["msg"] = err.Error()
		p.Abort("JSONError")
	}

	p.Ctx.Output.Status = http.StatusOK
	p.Data["json"] = device
	p.ServeJSON()
}
