package config

import (
	"werp/api/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectMySQL() *gorm.DB {
	// See https://github.com/go-sql-driver/mysql#dsn-data-source-name
	USER := "werp_beta_mysql_user"
	PASS := "root"
	PROTOCOL := "tcp(db:3306)"
	DATABASE := "werp_beta_general"
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
