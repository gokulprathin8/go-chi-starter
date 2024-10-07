package models

type User struct {
	ID       uint   `gorm:"primaryKey"` // Primary key for the User table
	Username string `gorm:"unique"`     // Unique username
	Password string // Hashed password
}
