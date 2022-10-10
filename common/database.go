package common

import (
	"fmt"
	"log"
	"ticketing/app"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var Db *gorm.DB

func init() {

	dsn := Configuration.DataBaseDsn

	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			NoLowerCase:   true,
		},
		Logger: logger.Discard,
	})

	if err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(&app.Ticket{}); err != nil {
		fmt.Print(err.Error())
	}

	if err := db.AutoMigrate(&app.Comment{}); err != nil {
		fmt.Print(err.Error())
	}

	if err := db.AutoMigrate(&app.React{}); err != nil {
		fmt.Print(err.Error())
	}

	Db = db
}
