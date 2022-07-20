package users_test

import (
	"testing"

	"github.com/reizt/ebra/conf"
)

func TestMain(m *testing.M) {
	conf.Migrate()
	m.Run()
}
