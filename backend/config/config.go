package config

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"ipicture-backend/models"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("ipicture.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败: ", err)
	}
	// 自动迁移表结构
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.File{})
} 