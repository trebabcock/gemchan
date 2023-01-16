package db

import (
	"log"

	"gemchan/app/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Init() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("/etc/gemchan/gorm.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	return DBMigrate(db)
}

func DBMigrate(db *gorm.DB) *gorm.DB {
	db.AutoMigrate(&model.Post{}, &model.Comment{}, &model.Board{})
	return db
}
