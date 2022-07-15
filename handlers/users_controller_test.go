package controllers_test

import (
	"ebra/config"
	"ebra/controllers"
	"ebra/models"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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
	config.Db = config.Db.Begin()
	config.Db.AutoMigrate(&models.User{})
	m.Run()
	config.Db.Rollback()
}

func TestGetAllUsersWhenUsersExist(t *testing.T) {
	// Given: Some users are registered
	users := []models.User{
		{Name: "John Smith"},
		{Name: "Michael Jackson"},
	}
	for i := 0; i < len(users); i++ {
		// Passing array and creating users at once will result all users' created_at are the same time.
		res := config.Db.Model(&models.User{}).Create(&users[i])
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

	// Then: Get JSON array having some user objects
	if assert.NoError(t, controllers.GetAllUsers(ctx)) {
		assert.Equal(t, rec.Code, http.StatusOK)
		assert.Equal(t, string(bytes)+"\n", rec.Body.String())
	}
	if err := config.Db.Delete(&users).Error; err != nil {
		panic(err)
	}
}

func TestGetAllUsersWhenUsersDontExist(t *testing.T) {
	// Given: No users are registered
	// When: GET /users
	ctx, _, rec := initTestContext(http.MethodGet, "/users", nil)

	// Then: Get JSON array having some user objects
	if assert.NoError(t, controllers.GetAllUsers(ctx)) {
		assert.Equal(t, rec.Code, http.StatusOK)
		assert.Equal(t, "[]"+"\n", rec.Body.String())
	}
}

func TestGetUserById(t *testing.T) {
	// Given: Some users are registered
	users := []models.User{
		{Name: "John Smith"},
		{Name: "Michael Jackson"},
	}
	for i := 0; i < len(users); i++ {
		// Passing array and creating users at once will result all users' created_at are the same time.
		res := config.Db.Model(&models.User{}).Create(&users[i])
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

	// Then: Get JSON array having some user objects
	if assert.NoError(t, controllers.GetUserById(ctx)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(bytes)+"\n", rec.Body.String())
	}
	if err := config.Db.Delete(&users).Error; err != nil {
		panic(err)
	}
}

func TestCreateUserWhenNameGiven(t *testing.T) {
	userJSON := `{"name": "John Smith"}`
	ctx, _, rec := initTestContext(http.MethodPost, "/users", strings.NewReader(userJSON))
	if assert.NoError(t, controllers.CreateUser(ctx)) {
		assert.Equal(t, rec.Code, http.StatusCreated)
		assert.Contains(t, rec.Body.String(), "John Smith")

		createdUser := &models.User{}
		json.Unmarshal(rec.Body.Bytes(), createdUser)
		assert.Equal(t, "John Smith", createdUser.Name)
	}
}
func TestCreateUserWhenNameNotGiven(t *testing.T) {
	userJSON := `{}`
	ctx, _, rec := initTestContext(http.MethodPost, "/users", strings.NewReader(userJSON))
	if assert.NoError(t, controllers.CreateUser(ctx)) {
		assert.Equal(t, rec.Code, http.StatusBadRequest)
	}
}
func TestCreateUserWhenBodyIsBlank(t *testing.T) {
	userJSON := ``
	ctx, _, rec := initTestContext(http.MethodPost, "/users", strings.NewReader(userJSON))
	if assert.NoError(t, controllers.CreateUser(ctx)) {
		assert.Equal(t, rec.Code, http.StatusBadRequest)
	}
}

func TestUpdateUserById(t *testing.T) {
	// Given: a user is registered
	user := &models.User{
		Name: "Robert Griesemer",
	}
	res := config.Db.Create(&user)
	if res.Error != nil {
		panic(res.Error)
	}

	// When: PATCH /users/:id with name
	userJSON := `{"name": "Ken Thompson"}`
	ctx, _, rec := initTestContext(http.MethodPatch, "/users/:id", strings.NewReader(userJSON))
	ctx.SetParamNames("id")
	ctx.SetParamValues(user.ID)

	// Then: Get status 200 and name changed
	if assert.NoError(t, controllers.UpdateUserById(ctx)) {
		assert.Equal(t, http.StatusNoContent, rec.Code)
		updatedUser := new(models.User)
		config.Db.First(&updatedUser, "id = ?", user.ID)
		assert.Equal(t, "Ken Thompson", updatedUser.Name) // Assert name changed successfully
	}
}
func TestUpdateUserByIdWhenBodyIsBlank(t *testing.T) {
	// Given: a user is registered
	user := &models.User{
		Name: "Robert Griesemer",
	}
	res := config.Db.Create(&user)
	if res.Error != nil {
		panic(res.Error)
	}

	// When: PATCH /users/:id with name
	userJSON := ``
	ctx, _, rec := initTestContext(http.MethodPatch, "/users/:id", strings.NewReader(userJSON))
	ctx.SetParamNames("id")
	ctx.SetParamValues(user.ID)

	// Then: Get status 200 and name changed
	if assert.NoError(t, controllers.UpdateUserById(ctx)) {
		assert.Equal(t, http.StatusNoContent, rec.Code)
		updatedUser := new(models.User)
		config.Db.First(&updatedUser, "id = ?", user.ID)
		assert.Equal(t, "Robert Griesemer", updatedUser.Name) // Assert name changed successfully
	}
}

func TestDeleteUserById(t *testing.T) {
	// Given: a user is registered
	user := &models.User{
		Name: "Robert Griesemer",
	}
	res := config.Db.Create(&user)
	if res.Error != nil {
		panic(res.Error)
	}

	// When: PATCH /users/:id with name
	ctx, _, rec := initTestContext(http.MethodDelete, "/users/:id", nil)
	ctx.SetParamNames("id")
	ctx.SetParamValues(user.ID)

	// Then: Get status 200 and name changed
	if assert.NoError(t, controllers.DeleteUserById(ctx)) {
		assert.Equal(t, http.StatusNoContent, rec.Code)
		// Assert user is deleted successfully
		deletedUser := new(models.User)
		res := config.Db.First(&deletedUser, "id = ?", user.ID)
		assert.EqualError(t, res.Error, "record not found")
	}
}
