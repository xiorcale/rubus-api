package controllers

import (
	"github.com/astaxie/beego"
)

// ErrorController serves JSON error to the client
type ErrorController struct {
	beego.Controller
}

// jsonError is the struct representing a JSON formatted error
type jsonError struct {
	Message string `json:"message"`
}

// ErrorJSONError serves an error with a JSON message
func (c *ErrorController) ErrorJSONError() {
	c.Ctx.Output.Status = c.Data["status"].(int)
	msg := c.Data["msg"].(string)
	c.Data["json"] = jsonError{Message: msg}
	c.ServeJSON()
}
