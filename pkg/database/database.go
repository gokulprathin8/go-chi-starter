package database

import (
	"chi-test/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Initialize() (*gorm.DB, error) {
	// Connect to SQLite database
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Automatically migrate the schema of the User struct
	db.AutoMigrate(&models.User{})
	return db, nil
}
