package config

import (
	"os"

	"github.com/reizt/ebra/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func loadEnv() {
	workdir := os.Getenv("WORKDIR")
	err := godotenv.Load(workdir + "/.env")
	if err != nil {
		panic(err)
	}
}

func ConnectMySQL() *gorm.DB {
	loadEnv()
	// See https://github.com/go-sql-driver/mysql
	type DsnConfig struct {
		Protocol string
		User     string
		Password string
		Database string
	}
	cnf := &DsnConfig{}
	switch os.Getenv("APP_ENV") {
	case "development":
		cnf = &DsnConfig{
			Protocol: "tcp(db:3306)",
			User:     os.Getenv("MYSQL_USER"),
			Password: os.Getenv("MYSQL_PASSWORD"),
			Database: os.Getenv("MYSQL_DATABASE"),
		}
	case "test":
		cnf = &DsnConfig{
			Protocol: "tcp(db:3306)",
			User:     "mysql_test_user",
			Password: "password",
			Database: "mysql_test_db",
		}
	}
	dsn := cnf.User + ":" + cnf.Password + "@" + cnf.Protocol + "/" + cnf.Database + "?charset=utf8mb4&parseTime=True&loc=Local"
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
