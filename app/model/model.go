package model

import (
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func DBMigrate(db *gorm.DB) *gorm.DB {
	db.AutoMigrate(&Post{}, &Comment{}, &Board{})
	return db
}
