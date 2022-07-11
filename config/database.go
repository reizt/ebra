package config

import (
	"ebra/models"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectMySQL() *gorm.DB {
	// See https://github.com/go-sql-driver/mysql#dsn-data-source-name
	err := godotenv.Load(os.Getenv("WORKDIR") + "/.env")
	if err != nil {
		panic(err)
	}
	USER := os.Getenv("MYSQL_USER")
	PASS := os.Getenv("MYSQL_PASSWORD")
	PROTOCOL := "tcp(db:3306)"
	DATABASE := os.Getenv("MYSQL_DATABASE")
	dsn := USER + ":" + PASS + "@" + PROTOCOL + "/" + DATABASE + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}
func Migrate() {
	db := ConnectMySQL()
	db.AutoMigrate(&models.User{})
}
