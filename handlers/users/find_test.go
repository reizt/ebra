package users_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/reizt/ebra/conf"
	handlers "github.com/reizt/ebra/handlers/users"
	"github.com/reizt/ebra/helpers"
	"github.com/reizt/ebra/models"
	"github.com/reizt/ebra/renderings"

	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	// Given: Some users are registered
	db := conf.ConnectMySQL()
	tx := db.Begin()
	users := []models.User{
		{Name: "John Smith"},
		{Name: "Michael Jackson"},
	}
	for i := 0; i < len(users); i++ {
		// Passing array and creating users at once will result all users' created_at are the same time.
		res := tx.Model(&models.User{}).Create(&users[i])
		if res.Error != nil {
			panic(res.Error)
		}
		time.Sleep(time.Millisecond * 100)
	}
	renderedUser := renderings.User{
		ID:        users[0].ID,
		Name:      users[0].Name,
		Email:     users[0].Email,
		CreatedAt: users[0].CreatedAt,
		UpdatedAt: users[0].UpdatedAt,
	}
	bytes, err := json.Marshal(renderedUser)
	if err != nil {
		panic(errors.New("failed to encode user to json"))
	}

	// When: GET /users/:id
	ctx, _, rec := helpers.InitTestContext(http.MethodGet, "/users/:id", nil)
	ctx.SetParamNames("id")
	ctx.SetParamValues(users[0].ID)
	ctx.Set(conf.DbContextKey, tx)

	// Then: Get JSON array having some user objects
	if assert.NoError(t, handlers.GetUserById(ctx)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(bytes)+"\n", rec.Body.String())
	}
	if err := tx.Delete(&users).Error; err != nil {
		panic(err)
	}
	tx.Rollback()
}
