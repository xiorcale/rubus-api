package controllers

import (
	"net/http"

	"github.com/astaxie/beego"
	"github.com/kjuvi/rubus-api/models"
)

// ErrorController serves JSON error to the client
type ErrorController struct {
	beego.Controller
}

// ErrorJSONError serves an error with a JSON message
func (c *ErrorController) ErrorJSONError() {
	err := c.Data["error"].(*models.JSONError)
	c.Ctx.Output.Status = err.Status
	c.Data["json"] = err
	c.ServeJSON()
}

// ErrorBadRequest serves a 400 Bad Request Error
func (c *ErrorController) ErrorBadRequest() {
	c.Ctx.Output.Status = http.StatusBadRequest
	err := c.Data["error"].(string)
	c.Data["json"] = map[string]string{"error": err}
	c.ServeJSON()
}
