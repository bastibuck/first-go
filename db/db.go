package db

import (
	"first-go/db/entities"
	"log"

	_ "github.com/mattn/go-sqlite3" // needed for the side-effects of this package (registering sqlite3 with database/sql)
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func OpenConnection() *gorm.DB {
	// Open a database connection
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&entities.Users{}, &entities.Events{})

	log.Println("Database connection established and tables prepared")

	return db
}
