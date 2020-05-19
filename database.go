package main

import (
	pg "github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/xiorcale/rubus-api/models"
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
