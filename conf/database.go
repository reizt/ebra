package conf

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/reizt/ebra/models"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectMySQL() *gorm.DB {
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
	dsn := fmt.Sprintf(
		"%s:%s@%s/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cnf.User, cnf.Password, cnf.Protocol, cnf.Database,
	)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Warn, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
	}), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}
	return db
}
func ConnectSessionDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("sessions." + os.Getenv("APP_ENV") + ".db"))
	if err != nil {
		panic(err)
	}
	return db
}
func Migrate() {
	db := ConnectMySQL()
	db.AutoMigrate(&models.User{})
	db = ConnectSessionDB()
	db.AutoMigrate(&models.Session{})
}
