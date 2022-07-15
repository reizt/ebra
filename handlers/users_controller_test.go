package handlers_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/reizt/ebra/config"
	"github.com/reizt/ebra/handlers"
	"github.com/reizt/ebra/models"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func initTestContext(method string, path string, body io.Reader) (ctx echo.Context, req *http.Request, rec *httptest.ResponseRecorder) {
	req = httptest.NewRequest(http.MethodGet, "/", body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	ctx = echo.New().NewContext(req, rec)
	ctx.SetPath("/users/:id")
	return ctx, req, rec
}

func TestMain(m *testing.M) {
	db := config.ConnectMySQL()
	db.AutoMigrate(&models.User{})
	m.Run()
}

func TestGetAllUsersWhenUsersExist(t *testing.T) {
	// Given: Some users are registered
	db := config.ConnectMySQL()
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
	// Reverse the users because GetAllUsers get users order by created_at desc
	for i := 0; i < len(users)/2; i++ {
		users[i], users[len(users)-i-1] = users[len(users)-i-1], users[i]
	}
	bytes, err := json.Marshal(users)
	if err != nil {
		panic(err)
	}

	// When: GET /users
	ctx, _, rec := initTestContext(http.MethodGet, "/users", nil)
	ctx.Set(config.DbContextKey, tx)

	// Then: Get JSON array having some user objects
	if assert.NoError(t, handlers.GetAllUsers(ctx)) {
		assert.Equal(t, rec.Code, http.StatusOK)
		assert.Equal(t, string(bytes)+"\n", rec.Body.String())
	}
	if err := db.Delete(&users).Error; err != nil {
		panic(err)
	}
	tx.Rollback()
}

func TestGetAllUsersWhenUsersDontExist(t *testing.T) {
	// Given: No users are registered
	db := config.ConnectMySQL()
	tx := db.Begin()
	// When: GET /users
	ctx, _, rec := initTestContext(http.MethodGet, "/users", nil)
	ctx.Set(config.DbContextKey, tx)

	// Then: Get JSON array having some user objects
	if assert.NoError(t, handlers.GetAllUsers(ctx)) {
		assert.Equal(t, rec.Code, http.StatusOK)
		assert.Equal(t, "[]"+"\n", rec.Body.String())
	}
	tx.Rollback()
}

func TestGetUserById(t *testing.T) {
	// Given: Some users are registered
	db := config.ConnectMySQL()
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
	bytes, err := json.Marshal(users[0])
	if err != nil {
		panic(errors.New("failed to encode user to json"))
	}

	// When: GET /users/:id
	ctx, _, rec := initTestContext(http.MethodGet, "/users/:id", nil)
	ctx.SetParamNames("id")
	ctx.SetParamValues(users[0].ID)
	ctx.Set(config.DbContextKey, tx)

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

func TestCreateUserWhenNameGiven(t *testing.T) {
	// Given: no users are registered
	db := config.ConnectMySQL()
	tx := db.Begin()
	userJSON := `{"name": "John Smith"}`

	// When: POST /users with body including property `name`
	ctx, _, rec := initTestContext(http.MethodPost, "/users", strings.NewReader(userJSON))
	ctx.Set(config.DbContextKey, tx)

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
	db := config.ConnectMySQL()
	tx := db.Begin()

	ctx, _, rec := initTestContext(http.MethodPost, "/users", strings.NewReader(userJSON))
	ctx.Set(config.DbContextKey, tx)

	// Then: Can't create user
	if assert.NoError(t, handlers.CreateUser(ctx)) {
		assert.Equal(t, rec.Code, http.StatusBadRequest)
	}
	tx.Rollback()
}
func TestCreateUserWhenBodyIsBlank(t *testing.T) {
	db := config.ConnectMySQL()
	tx := db.Begin()
	// When: POST /users without body
	ctx, _, rec := initTestContext(http.MethodPost, "/users", nil)
	ctx.Set(config.DbContextKey, tx)
	// Given: Can't create user
	if assert.NoError(t, handlers.CreateUser(ctx)) {
		assert.Equal(t, rec.Code, http.StatusBadRequest)
	}
	tx.Rollback()
}

func TestUpdateUserById(t *testing.T) {
	// Given: a user is registered
	db := config.ConnectMySQL()
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
	ctx, _, rec := initTestContext(http.MethodPatch, "/users/:id", strings.NewReader(userJSON))
	ctx.SetParamNames("id")
	ctx.SetParamValues(user.ID)
	ctx.Set(config.DbContextKey, tx)

	// Then: Get status 200 and name changed
	if assert.NoError(t, handlers.UpdateUserById(ctx)) {
		assert.Equal(t, http.StatusNoContent, rec.Code)
		updatedUser := new(models.User)
		tx.First(&updatedUser, "id = ?", user.ID)
		assert.Equal(t, "Ken Thompson", updatedUser.Name) // Assert name changed successfully
	}
	tx.Rollback()
}
func TestUpdateUserByIdWhenBodyIsBlank(t *testing.T) {
	// Given: a user is registered
	user := &models.User{
		Name: "Robert Griesemer",
	}
	db := config.ConnectMySQL()
	tx := db.Begin()
	res := tx.Create(&user)
	if res.Error != nil {
		panic(res.Error)
	}

	// When: PATCH /users/:id with name
	userJSON := ``
	ctx, _, rec := initTestContext(http.MethodPatch, "/users/:id", strings.NewReader(userJSON))
	ctx.SetParamNames("id")
	ctx.SetParamValues(user.ID)
	ctx.Set(config.DbContextKey, tx)

	// Then: Get status 200 and name changed
	if assert.NoError(t, handlers.UpdateUserById(ctx)) {
		assert.Equal(t, http.StatusNoContent, rec.Code)
		updatedUser := new(models.User)
		tx.First(&updatedUser, "id = ?", user.ID)
		assert.Equal(t, "Robert Griesemer", updatedUser.Name) // Assert name changed successfully
	}
	tx.Rollback()
}

func TestDeleteUserById(t *testing.T) {
	// Given: a user is registered
	db := config.ConnectMySQL()
	tx := db.Begin()
	user := &models.User{
		Name: "Robert Griesemer",
	}
	res := tx.Create(&user)
	if res.Error != nil {
		panic(res.Error)
	}

	// When: PATCH /users/:id with name
	ctx, _, rec := initTestContext(http.MethodDelete, "/users/:id", nil)
	ctx.SetParamNames("id")
	ctx.SetParamValues(user.ID)
	ctx.Set(config.DbContextKey, tx)

	// Then: Get status 200 and name changed
	if assert.NoError(t, handlers.DeleteUserById(ctx)) {
		assert.Equal(t, http.StatusNoContent, rec.Code)
		// Assert user is deleted successfully
		deletedUser := new(models.User)
		res := tx.First(&deletedUser, "id = ?", user.ID)
		assert.EqualError(t, res.Error, "record not found")
	}
	tx.Rollback()
}
