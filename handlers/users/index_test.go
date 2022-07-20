package users_test

import (
	"encoding/json"
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

func TestGetUsersWhenUsersExist(t *testing.T) {
	// Given: Some users are registered
	db := conf.ConnectMySQL()
	tx := db.Begin()
	users := []models.User{
		{Name: "John Smith"},
		{Name: "Michael Jackson"},
	}
	for i := 0; i < len(users); i++ {
		// Passing array and creating users at once will result all users' created_at are the same time.
		res := db.Model(&models.User{}).Create(&users[i])
		if res.Error != nil {
			panic(res.Error)
		}
		time.Sleep(time.Millisecond * 100)
	}
	// Reverse the users because GetUsers get users order by created_at desc
	for i := 0; i < len(users)/2; i++ {
		users[i], users[len(users)-i-1] = users[len(users)-i-1], users[i]
	}
	userRenderings := []renderings.User{}
	for _, u := range users {
		userRenderings = append(userRenderings, renderings.User{
			ID:        u.ID,
			Name:      u.Name,
			Email:     u.Email,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		})
	}
	bytes, err := json.Marshal(userRenderings)
	if err != nil {
		panic(err)
	}

	// When: GET /users
	ctx, _, rec := helpers.InitTestContext(http.MethodGet, "/users", nil)
	ctx.Set(conf.DbContextKey, tx)

	// Then: Get JSON array having some user objects
	if assert.NoError(t, handlers.GetUsers(ctx)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(bytes)+"\n", rec.Body.String())
	}
	if err := db.Delete(&users).Error; err != nil {
		panic(err)
	}
	tx.Rollback()
}

func TestGetUsersWhenUsersDontExist(t *testing.T) {
	// Given: No users are registered
	db := conf.ConnectMySQL()
	tx := db.Begin()
	// When: GET /users
	ctx, _, rec := helpers.InitTestContext(http.MethodGet, "/users", nil)
	ctx.Set(conf.DbContextKey, tx)

	// Then: Get JSON array having some user objects
	if assert.NoError(t, handlers.GetUsers(ctx)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "[]"+"\n", rec.Body.String())
	}
	tx.Rollback()
}
