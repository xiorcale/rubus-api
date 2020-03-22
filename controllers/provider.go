package controllers

import (
	"encoding/json"

	"net/http"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/kjuvi/rubus-api/models"
)

// Operations offered by the provider
type ProviderController struct {
	beego.Controller
}

// GetDevice returns a `Device`
// @Title Get
// @Description Get the Rubus `Device` with the given `deviceId`
// @Param	deviceId		path 	int	true		"The device id to get"
// @Success 200 {object} models.Device
// @router /device/:deviceId [get]
func (p *ProviderController) GetDevice() {
	port := p.GetString(":deviceId")
	url := "http://rubus_provider:1080/device/" + port

	device := models.Device{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		p.Data["status"] = http.StatusInternalServerError
		p.Data["msg"] = "Internal Server Error"
		p.Abort("JSONError")
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		p.Data["status"] = http.StatusInternalServerError
		p.Data["msg"] = "Internal Server Error"
		p.Abort("JSONError")
	}

	if err := device.Bind(res.Body); err != nil {
		p.Data["status"] = http.StatusInternalServerError
		p.Data["msg"] = "Internal Server Error"
		p.Abort("JSONError")
	}
	defer res.Body.Close()

	p.Ctx.Output.Status = http.StatusOK
	p.Data["json"] = device
	p.ServeJSON()
}

// GetAllDevices requests the provider to send the list of available `Device`
// for Rubus to provision
// @Title GetAll
// @Description get all the rubus `Device`
// @Success 200 {object} []models.Device
// @router /device [get]
func (p *ProviderController) GetAllDevices() {
	url := "http://rubus_provider:1080/device"
	devices := []models.Device{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		p.Data["status"] = http.StatusInternalServerError
		p.Data["msg"] = "Internal Server Error"
		p.Abort("JSONError")
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		p.Data["status"] = http.StatusInternalServerError
		p.Data["msg"] = "Internal Server Error"
		p.Abort("JSONError")
	}

	// buf := new(bytes.Buffer)
    // buf.ReadFrom(res.Body)
	// newStr := buf.String()
	
	// logs.Debug(newStr)

	if err := json.NewDecoder(res.Body).Decode(&devices); err != nil {
		logs.Debug(err)
		p.Data["status"] = http.StatusInternalServerError
		p.Data["msg"] = "Internal Server Error"
		p.Abort("JSONError")
	}
	defer res.Body.Close()

	p.Ctx.Output.Status = http.StatusOK
	p.Data["json"] = devices
	p.ServeJSON()
}

// PowerOn boots the `Device` on the given `port`
// @Title PowerOn
// @Description boots the `Device` on the given `port`
// @Success 204
// @router /device/:deviceId/on [post]
func (p *ProviderController) PowerOn() {
	port := p.GetString(":deviceId")
	url := "http://rubus_provider:1080/device/" + port + "/on"

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		p.Data["status"] = http.StatusInternalServerError
		p.Data["msg"] = "Internal Server Error"
		p.Abort("JSONError")
	}

	if _, err = http.DefaultClient.Do(req); err != nil {
		p.Data["status"] = http.StatusInternalServerError
		p.Data["msg"] = "Internal Server Error"
		p.Abort("JSONError")
	}

	p.Ctx.Output.Status = http.StatusNoContent
	p.ServeJSON()
}

// PowerOff shuts down the `Device` on the given `port`
// @Title PowerOff
// @Description shuts down the `Device` on the given `port`
// @Success 204
// @router /device/:deviceId/off [post]
func (p *ProviderController) PowerOff() {
	port := p.GetString(":deviceId")
	url := "http://rubus_provider:1080/device/" + port + "/off"

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		p.Data["status"] = http.StatusInternalServerError
		p.Data["msg"] = "Internal Server Error"
		p.Abort("JSONError")
	}

	if _, err = http.DefaultClient.Do(req); err != nil {
		p.Data["status"] = http.StatusInternalServerError
		p.Data["msg"] = "Internal Server Error"
		p.Abort("JSONError")
	}

	p.Ctx.Output.Status = http.StatusNoContent
	p.ServeJSON()
}
