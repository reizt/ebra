package users_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/reizt/ebra/conf"
	handlers "github.com/reizt/ebra/handlers/users"
	"github.com/reizt/ebra/models"

	"github.com/stretchr/testify/assert"
)

func TestUpdateUser(t *testing.T) {
	// Given: a user is registered
	db := conf.ConnectMySQL()
	tx := db.Begin()
	user := &models.User{
		Name: "Robert Griesemer",
	}
	res := tx.Create(&user)
	if res.Error != nil {
		panic(res.Error)
	}

	// When: PATCH /users/:id with name
	userJSON := `{"name": "Ken Thompson"}`
	ctx, _, rec := InitTestContext(http.MethodPatch, "/users/:id", strings.NewReader(userJSON))
	ctx.SetParamNames("id")
	ctx.SetParamValues(user.ID)
	ctx.Set(conf.DbContextKey, tx)

	// Then: Get status 200 and name changed
	if assert.NoError(t, handlers.UpdateUser(ctx)) {
		assert.Equal(t, http.StatusNoContent, rec.Code)
		updatedUser := new(models.User)
		tx.First(&updatedUser, "id = ?", user.ID)
		assert.Equal(t, "Ken Thompson", updatedUser.Name) // Assert name changed successfully
	}
	tx.Rollback()
}
func TestUpdateUserWhenBodyIsBlank(t *testing.T) {
	// Given: a user is registered
	user := &models.User{
		Name: "Robert Griesemer",
	}
	db := conf.ConnectMySQL()
	tx := db.Begin()
	res := tx.Create(&user)
	if res.Error != nil {
		panic(res.Error)
	}

	// When: PATCH /users/:id with name
	userJSON := ``
	ctx, _, rec := InitTestContext(http.MethodPatch, "/users/:id", strings.NewReader(userJSON))
	ctx.SetParamNames("id")
	ctx.SetParamValues(user.ID)
	ctx.Set(conf.DbContextKey, tx)

	// Then: Get status 200 and name changed
	if assert.NoError(t, handlers.UpdateUser(ctx)) {
		assert.Equal(t, http.StatusNoContent, rec.Code)
		updatedUser := new(models.User)
		tx.First(&updatedUser, "id = ?", user.ID)
		assert.Equal(t, "Robert Griesemer", updatedUser.Name) // Assert name changed successfully
	}
	tx.Rollback()
}
