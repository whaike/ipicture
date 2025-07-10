package main

import (
	"ipicture-backend/api"
	"ipicture-backend/config"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func main() {
	config.InitDB()
	r := gin.Default()
	r.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"*"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
}))

	r.Static("/thumbnails", "./thumbnails")
	// 新增：暴露原图目录
	r.Static("/uploads", "./uploads")
	api.RegisterRoutes(r)
	api.RegisterFileRoutes(r)

	r.Run(":8080")
} 