package database

import (
	"ambassador/src/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	var err error

	DB, err = gorm.Open(mysql.Open("root:root@tcp(db:3306)/ambassador"), &gorm.Config{})

	if err != nil {
		panic("could not connect to db")
	}
}

func AutoMigrate() {
	DB.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Link{},
		&models.Order{},
		&models.OrderItem{},
	)
}
