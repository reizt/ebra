package auth_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/reizt/ebra/middlewares"

	"github.com/labstack/echo/v4"

	"github.com/reizt/ebra/conf"
	handlers "github.com/reizt/ebra/handlers/auth"
	"github.com/reizt/ebra/helpers"
	"github.com/reizt/ebra/models"
	"github.com/stretchr/testify/assert"
)

func TestSigninWithCorrectPassword(t *testing.T) {
	db := conf.ConnectMySQL()
	tx := db.Begin()
	user := new(models.User)
	user.Name = "John Smith"
	user.Email = "john.smith@example.com"
	user.Password = "password"
	if err := tx.Create(&user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	signinReqJSON := `{"email": "john.smith@example.com", "password": "password"}`
	ctx, _, rec := helpers.InitTestContext(http.MethodPost, "/auth/signin", strings.NewReader(signinReqJSON))
	ctx.Set(conf.DbContextKey, tx)

	if assert.NoError(t, handlers.Signin(ctx)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		jwt := rec.Header().Get(echo.HeaderSetCookie)
		assert.NotNil(t, jwt)
		claims, err := middlewares.ValidateJwt(jwt)
		if err != nil {
			t.Error(err)
		}
		assert.NotNil(t, claims["sessionId"])
		assert.Contains(t, rec.Body.String(), "John Smith")
		assert.Contains(t, rec.Body.String(), "john.smith@example.com")
		assert.NotContains(t, rec.Body.String(), "PasswordDigest")
	}
	tx.Rollback()
}
func TestSigninWithIncorrectPassword(t *testing.T) {
	db := conf.ConnectMySQL()
	tx := db.Begin()
	user := new(models.User)
	user.Name = "John Smith"
	user.Email = "john.smith@example.com"
	user.Password = "correct"
	if err := tx.Create(&user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	signinReqJSON := `{"email": "john.smith@example.com", "password": "incorrect"}`
	ctx, _, rec := helpers.InitTestContext(http.MethodPost, "/auth/signin", strings.NewReader(signinReqJSON))
	ctx.Set(conf.DbContextKey, tx)

	if assert.NoError(t, handlers.Signin(ctx)) {
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		jwt := rec.Header().Get(echo.HeaderSetCookie)
		assert.Equal(t, "", jwt)
	}
	tx.Rollback()
}
