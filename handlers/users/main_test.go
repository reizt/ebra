package users_test

import (
	"testing"

	"github.com/reizt/ebra/conf"
	"github.com/reizt/ebra/models"
)

func TestMain(m *testing.M) {
	db := conf.ConnectMySQL()
	db.AutoMigrate(&models.User{})
	m.Run()
}
