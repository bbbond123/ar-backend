package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// 从环境变量读取数据库配置
	host := os.Getenv("BLUEPRINT_DB_HOST")
	user := os.Getenv("BLUEPRINT_DB_USERNAME")
	password := os.Getenv("BLUEPRINT_DB_PASSWORD")
	dbname := os.Getenv("BLUEPRINT_DB_DATABASE")
	port := os.Getenv("BLUEPRINT_DB_PORT")

	// 验证必需的环境变量
	if host == "" || user == "" || password == "" || dbname == "" || port == "" {
		log.Fatal("Database configuration is incomplete. Please check these environment variables: BLUEPRINT_DB_HOST, BLUEPRINT_DB_USERNAME, BLUEPRINT_DB_PASSWORD, BLUEPRINT_DB_DATABASE, BLUEPRINT_DB_PORT")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tokyo",
		host, user, password, dbname, port)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database!", err)
	}
	DB = database
}

func GetDB() *gorm.DB {
	return DB
}
