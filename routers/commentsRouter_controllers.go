package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/kjuvi/rubus-api/controllers:ProviderController"] = append(beego.GlobalControllerRouter["github.com/kjuvi/rubus-api/controllers:ProviderController"],
        beego.ControllerComments{
            Method: "GetAllDevices",
            Router: `/device`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/kjuvi/rubus-api/controllers:ProviderController"] = append(beego.GlobalControllerRouter["github.com/kjuvi/rubus-api/controllers:ProviderController"],
        beego.ControllerComments{
            Method: "GetDevice",
            Router: `/device/:deviceId`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/kjuvi/rubus-api/controllers:ProviderController"] = append(beego.GlobalControllerRouter["github.com/kjuvi/rubus-api/controllers:ProviderController"],
        beego.ControllerComments{
            Method: "PowerOff",
            Router: `/device/:deviceId/off`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/kjuvi/rubus-api/controllers:ProviderController"] = append(beego.GlobalControllerRouter["github.com/kjuvi/rubus-api/controllers:ProviderController"],
        beego.ControllerComments{
            Method: "PowerOn",
            Router: `/device/:deviceId/on`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/kjuvi/rubus-api/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/kjuvi/rubus-api/controllers:UserController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/kjuvi/rubus-api/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/kjuvi/rubus-api/controllers:UserController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/kjuvi/rubus-api/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/kjuvi/rubus-api/controllers:UserController"],
        beego.ControllerComments{
            Method: "Get",
            Router: `/:uid`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/kjuvi/rubus-api/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/kjuvi/rubus-api/controllers:UserController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/:uid`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/kjuvi/rubus-api/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/kjuvi/rubus-api/controllers:UserController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/:uid`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/kjuvi/rubus-api/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/kjuvi/rubus-api/controllers:UserController"],
        beego.ControllerComments{
            Method: "Login",
            Router: `/login`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
