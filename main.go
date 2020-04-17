package main

import (
	"fmt"

	"github.com/kjuvi/rubus-api/controllers"
	"github.com/kjuvi/rubus-api/models"
	_ "github.com/kjuvi/rubus-api/routers"
	"golang.org/x/crypto/bcrypt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/lib/pq"
)

func init() {
	orm.RegisterDriver("postgres", orm.DRPostgres)

	user := beego.AppConfig.String("user")

	password := beego.AppConfig.String("password")
	host := beego.AppConfig.String("host")
	dbName := beego.AppConfig.String("dbname")
	dataSource := fmt.Sprintf(
		"user=%s password=%s host=%s dbname=%s sslmode=disable",
		user, password, host, dbName,
	)
	orm.RegisterDataBase(
		"default",
		"postgres",
		dataSource,
	)
}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	name := "default"
	force := true
	verbose := true

	if err := orm.RunSyncdb(name, force, verbose); err != nil {
		panic(err)
	}

	// TODO: change this for prod
	// insert a default administrator
	cost, _ := beego.AppConfig.Int("hashcost")
	bytes, _ := bcrypt.GenerateFromPassword([]byte("rubus_secret"), cost)
	admin := models.User{
		Username:     "admin",
		Email:        "admin@mail.com",
		PasswordHash: string(bytes),
		Role:         models.EnumRoleAdmin,
	}
	models.AddUser(&admin)

	beego.ErrorController(&controllers.ErrorController{})

	beego.Run()
}
