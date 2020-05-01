// @APIVersion 1.0.0
// @Title Rubus API
// @Description Rubus API exposes provisioning services to manage an edge cluster system (i.e. Raspberry pi). This API takes advantage of various HTTP features like authentication, verbs or status code. All requests and response bodies are JSON encoded, including error responses.
// @Contact quentin.vaucher@master.hes-so.ch
// @License MIT
// @LicenseUrl https://opensource.org/licenses/MIT
package routers

import (
	"github.com/kjuvi/rubus-api/controllers"
	"github.com/kjuvi/rubus-api/services"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/admin",
			beego.NSInclude(
				&controllers.AdminController{},
			),
		),
		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
		beego.NSNamespace("/device",
			beego.NSInclude(
				&controllers.DeviceController{},
				&controllers.ProvisionerController{},
			),
		),
	)
	beego.AddNamespace(ns)

	// authentication middleware
	beego.InsertFilter("/*", beego.BeforeRouter, services.FilterUser)
}
