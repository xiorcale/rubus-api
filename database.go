package main

import (
	pg "github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/labstack/echo"
	"github.com/xiorcale/rubus-api/models"
	"golang.org/x/crypto/bcrypt"
)

var modelsList = []interface{}{
	(*models.User)(nil),
	(*models.Device)(nil),
}

func createSchema(db *pg.DB) error {
	for _, model := range modelsList {
		if err := db.CreateTable(model, &orm.CreateTableOptions{IfNotExists: true}); err != nil {
			return err
		}
	}

	return nil
}

func deleteSchema(db *pg.DB) error {
	for _, model := range modelsList {
		if err := db.DropTable(model, &orm.DropTableOptions{IfExists: true, Cascade: true}); err != nil {
			return err
		}
	}

	return nil
}

func createAdmin(s server) error {
	dbCfg := s.cfg.Section("database")
	admin := dbCfg.Key("admin_username").String()
	email := dbCfg.Key("admin_email").String()
	password := dbCfg.Key("admin_password").String()

	cost, _ := s.cfg.Section("security").Key("hashcost").Int()
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), cost)

	user := models.User{
		Username:     admin,
		Email:        email,
		PasswordHash: string(bytes),
		Role:         models.EnumRoleAdmin,
	}

	if jsonErr := models.AddUser(s.db, &user); jsonErr != nil {
		return echo.NewHTTPError(jsonErr.Status, jsonErr)
	}

	return nil
}
