package main

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/xiorcale/rubus-api/controllers"
	_ "github.com/xiorcale/rubus-api/docs"
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
    authentication := controllers.AuthenticationController{DB: s.db, Cfg: s.cfg}
	user := controllers.UserController{DB: s.db, Cfg: s.cfg}
	device := controllers.DeviceController{DB: s.db}
	provisioner := controllers.ProvisionerController{DB: s.db}
	admin := controllers.AdminController{DB: s.db, Cfg: s.cfg}

	// groups
	userGr := s.e.Group("/user")
	deviceGr := s.e.Group("/device")
	adminGr := s.e.Group("/admin")

	// jwt protection
	secret := s.cfg.Section("security").Key("jwtsecret").String()
	userGr.Use(middleware.JWT([]byte(secret)))
	deviceGr.Use(middleware.JWT([]byte(secret)))
	adminGr.Use(middleware.JWT([]byte(secret)))

	s.e.GET("/login", authentication.Login)

	// user endpoints
	userGr.GET("/me", user.GetMe)
	userGr.PUT("/me", user.UpdateMe)
	userGr.DELETE("/me", user.DeleteMe)

	// device endpoints
	deviceGr.GET("", device.ListDevice)
	deviceGr.GET("/:id", device.Get)
	deviceGr.POST("/:id/on", device.PowerOn)
	deviceGr.POST("/:id/off", device.PowerOff)
	deviceGr.POST("/:id/acquire", provisioner.Acquire)
	deviceGr.POST("/:id/release", provisioner.Release)
	deviceGr.POST("/:id/deploy", provisioner.Deploy)

	// admin endpoints
	adminGr.POST("/device", admin.CreateDevice)
	adminGr.DELETE("/device", admin.DeleteDevice)
	adminGr.POST("/user", admin.CreateUser)
	adminGr.GET("/user", admin.ListUser)
	adminGr.DELETE("/user/:id", admin.DeleteUser)
	adminGr.POST("/user/:id/expiration", admin.UpdateUserExpiration)
}
