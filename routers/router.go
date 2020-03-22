// @APIVersion 1.0.0
// @Title Rubus API
// @Description Rubus API exposes provisioning services to manage an edge cluster system (i.e. Raspberry pi). This API takes advantage of various HTTP features like authentication, verbs or status code. All requests and response bodies are JSON encoded, including error responses.
// @Contact quentin.vaucher@master.hes-so.ch
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/kjuvi/rubus-api/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
		beego.NSNamespace("/provider",
			beego.NSInclude(
				&controllers.ProviderController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
