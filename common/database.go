package common

import (
	"fmt"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"ticketing/app"
)

func OpenDb() *gorm.DB {
	dsn := Configuration.DataBaseDsn
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			NoLowerCase:   true,
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	return db
}

func init() {

	db := OpenDb()

	if err := db.AutoMigrate(&app.Ticket{}); err != nil {
		fmt.Print(err.Error())
	}

	if err := db.AutoMigrate(&app.Comment{}); err != nil {
		fmt.Print(err.Error())
	}

	if err := db.AutoMigrate(&app.React{}); err != nil {
		fmt.Print(err.Error())
	}
}
