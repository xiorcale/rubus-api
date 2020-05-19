package main

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/xiorcale/rubus-api/controllers"

	_ "github.com/xiorcale/docs"
)

func createRESTEndpoints(s server) {
	// middleware
	s.e.Use(middleware.Logger())
	s.e.Use(middleware.Recover())
	s.e.Use(middleware.GzipWithConfig((middleware.GzipConfig{
		Skipper: func(c echo.Context) bool {
			if strings.Contains(c.Request().URL.Path, "swagger") {
				return true
			}
			return false
		},
	})))
	s.e.Use(middleware.CORS())

	// documentation
	s.e.GET("/swagger/*", echoSwagger.WrapHandler)

	// controllers
	user := controllers.UserController{DB: s.db, Cfg: s.cfg}
	device := controllers.DeviceController{DB: s.db}
	provisioner := controllers.ProvisionerController{DB: s.db}
	admin := controllers.AdminController{DB: s.db, Cfg: s.cfg}

	// groups
	userGr := s.e.Group("/user")
	deviceGr := s.e.Group("/device")
	adminGr := s.e.Group("/admin")

	// jwt protection
	secret := s.cfg.Section("jwt").Key("jwtsecret").String()
	userGr.Use(middleware.JWT([]byte(secret)))
	deviceGr.Use(middleware.JWT([]byte(secret)))
	adminGr.Use(middleware.JWT([]byte(secret)))

	// user endpoints
	userGr.GET("", user.ListUser)
	userGr.GET("/login", user.Login)
	userGr.GET("/me", user.GetMe)
	userGr.PUT("/me", user.UpdateMe)
	userGr.DELETE("/me", user.DeleteMe)

	// device endpoints
	deviceGr.GET("", device.ListDevice)
	deviceGr.GET("/:deviceId", device.Get)
	deviceGr.POST("/:deviceId/on", device.PowerOn)
	deviceGr.POST("/:deviceId/off", device.PowerOff)
	deviceGr.POST("/:deviceId/acquire", provisioner.Acquire)
	deviceGr.POST("/:deviceId/release", provisioner.Release)
	deviceGr.POST("/:deviceId/deploy", provisioner.Deploy)

	// admin endpoints
	adminGr.POST("/device", admin.CreateDevice)
	adminGr.DELETE("/device", admin.DeleteDevice)
	adminGr.POST("/user", admin.CreateUser)
}
