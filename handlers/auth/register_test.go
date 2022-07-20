package auth_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/reizt/ebra/conf"
	handlers "github.com/reizt/ebra/handlers/auth"
	"github.com/reizt/ebra/helpers"
	"github.com/reizt/ebra/models"

	"github.com/stretchr/testify/assert"
)

func TestRegisterWhenParametersGiven(t *testing.T) {
	// Given: no users are registered
	db := conf.ConnectMySQL()
	tx := db.Begin()
	userJSON := `{"name": "John Smith", "email": "john.smith@example.com", "password": "password"}`

	// When: POST /users with body including property `name`
	ctx, _, rec := helpers.InitTestContext(http.MethodPost, "/auth/register", strings.NewReader(userJSON))
	ctx.Set(conf.DbContextKey, tx)

	// Then: Successfully user is created
	if assert.NoError(t, handlers.Register(ctx)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Contains(t, rec.Body.String(), "John Smith")
		assert.Contains(t, rec.Body.String(), "john.smith@example.com")
		assert.NotContains(t, rec.Body.String(), "PasswordDigest")

		createdUser := &models.User{}
		json.Unmarshal(rec.Body.Bytes(), createdUser)
		assert.Equal(t, "John Smith", createdUser.Name)
	}
	tx.Rollback()
}
func TestRegisterWhenParametersNotGiven(t *testing.T) {
	// When: POST /users without property `name`
	userJSON := `{}`
	db := conf.ConnectMySQL()
	tx := db.Begin()

	ctx, _, rec := helpers.InitTestContext(http.MethodPost, "/auth/register", strings.NewReader(userJSON))
	ctx.Set(conf.DbContextKey, tx)

	// Then: Can't create user
	if assert.NoError(t, handlers.Register(ctx)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
	tx.Rollback()
}
func TestRegisterWhenBodyIsBlank(t *testing.T) {
	db := conf.ConnectMySQL()
	tx := db.Begin()
	// When: POST /users without body
	ctx, _, rec := helpers.InitTestContext(http.MethodPost, "/auth/register", nil)
	ctx.Set(conf.DbContextKey, tx)
	// Given: Can't create user
	if assert.NoError(t, handlers.Register(ctx)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
	tx.Rollback()
}
