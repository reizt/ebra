package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/reizt/ebra/bindings"
	"github.com/reizt/ebra/conf"
	"github.com/reizt/ebra/middlewares"
	"github.com/reizt/ebra/models"
)

func Signin(c echo.Context) error {
	db := c.Get(conf.DbContextKey).(*gorm.DB)
	params := &bindings.SigninParams{}
	currentUser := &models.User{}
	if err := (&echo.DefaultBinder{}).BindBody(c, &params); err != nil {
		return err
	}
	res := db.First(&currentUser, "email = ?", params.Email)
	if res.Error != nil {
		return echo.ErrUnauthorized
	}
	if err := bcrypt.CompareHashAndPassword([]byte(currentUser.PasswordDigest), []byte(params.Password)); err != nil {
		return echo.ErrUnauthorized
	}
	sessionId, err := middlewares.StartSession(currentUser.ID)
	if err != nil {
		return err
	}
	cookie, err := middlewares.GenerateJwt(sessionId)
	if err != nil {
		return err
	}
	c.Response().Header().Set(echo.HeaderSetCookie, cookie)
	return c.JSON(http.StatusOK, currentUser)
}
