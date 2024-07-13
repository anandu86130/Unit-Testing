package database

import (
	"log"
	"userPage/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func autoMigrate(db *gorm.DB) {
	db.AutoMigrate(&models.Users{})
}

func CreateDB() {
	DSN := "host=localhost user=postgres password=rapunzel dbname=admin port=5432"
	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	autoMigrate(db)
	DB = db
}

func SetDB(database *gorm.DB) {
	DB = database
}
