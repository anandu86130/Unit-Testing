package models

import "gorm.io/gorm"

type Users struct {
	gorm.Model
	Name     string `gorm:"not null"`
	Email    string `gorm:"not null;unique"`
	Password string `gorm:"not null"`
}

type Admin struct {
	Name     string
	Email    string
	Password string
}
