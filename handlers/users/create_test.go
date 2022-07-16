package users_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/reizt/ebra/conf"
	handlers "github.com/reizt/ebra/handlers/users"
	"github.com/reizt/ebra/models"

	"github.com/stretchr/testify/assert"
)

func TestCreateUserWhenNameGiven(t *testing.T) {
	// Given: no users are registered
	db := conf.ConnectMySQL()
	tx := db.Begin()
	userJSON := `{"name": "John Smith"}`

	// When: POST /users with body including property `name`
	ctx, _, rec := InitTestContext(http.MethodPost, "/users", strings.NewReader(userJSON))
	ctx.Set(conf.DbContextKey, tx)

	// Then: Successfully user is created
	if assert.NoError(t, handlers.CreateUser(ctx)) {
		assert.Equal(t, rec.Code, http.StatusCreated)
		assert.Contains(t, rec.Body.String(), "John Smith")

		createdUser := &models.User{}
		json.Unmarshal(rec.Body.Bytes(), createdUser)
		assert.Equal(t, "John Smith", createdUser.Name)
	}
	tx.Rollback()
}
func TestCreateUserWhenNameNotGiven(t *testing.T) {
	// When: POST /users without property `name`
	userJSON := `{}`
	db := conf.ConnectMySQL()
	tx := db.Begin()

	ctx, _, rec := InitTestContext(http.MethodPost, "/users", strings.NewReader(userJSON))
	ctx.Set(conf.DbContextKey, tx)

	// Then: Can't create user
	if assert.NoError(t, handlers.CreateUser(ctx)) {
		assert.Equal(t, rec.Code, http.StatusBadRequest)
	}
	tx.Rollback()
}
func TestCreateUserWhenBodyIsBlank(t *testing.T) {
	db := conf.ConnectMySQL()
	tx := db.Begin()
	// When: POST /users without body
	ctx, _, rec := InitTestContext(http.MethodPost, "/users", nil)
	ctx.Set(conf.DbContextKey, tx)
	// Given: Can't create user
	if assert.NoError(t, handlers.CreateUser(ctx)) {
		assert.Equal(t, rec.Code, http.StatusBadRequest)
	}
	tx.Rollback()
}
