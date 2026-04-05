package db

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func GetGormDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("file:data.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("打开数据库失败:", err)
	}
	return db, nil
}
