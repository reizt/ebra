package users_test

import (
	"net/http"
	"testing"

	"github.com/reizt/ebra/conf"
	handlers "github.com/reizt/ebra/handlers/users"
	"github.com/reizt/ebra/models"

	"github.com/stretchr/testify/assert"
)

func TestDeleteUser(t *testing.T) {
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
	ctx, _, rec := InitTestContext(http.MethodDelete, "/users/:id", nil)
	ctx.SetParamNames("id")
	ctx.SetParamValues(user.ID)
	ctx.Set(conf.DbContextKey, tx)

	// Then: Get status 200 and name changed
	if assert.NoError(t, handlers.DeleteUser(ctx)) {
		assert.Equal(t, http.StatusNoContent, rec.Code)
		// Assert user is deleted successfully
		deletedUser := new(models.User)
		res := tx.First(&deletedUser, "id = ?", user.ID)
		assert.EqualError(t, res.Error, "record not found")
	}
	tx.Rollback()
}
