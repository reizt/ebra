package auth

import (
	"net/http"

	"github.com/reizt/ebra/renderings"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/reizt/ebra/bindings"
	"github.com/reizt/ebra/conf"
	"github.com/reizt/ebra/middlewares"
	"github.com/reizt/ebra/models"
)

func Signin(c echo.Context) (err error) {
	db := c.Get(conf.DbContextKey).(*gorm.DB)
	params := &bindings.SigninParams{}
	resp := &renderings.UserResponse{}
	currentUser := &models.User{}
	err = (&echo.DefaultBinder{}).BindBody(c, &params)
	if err != nil || params.Email == "" || params.Password == "" {
		return c.NoContent(http.StatusBadRequest)
	}
	res := db.First(&currentUser, "email = ?", params.Email)
	if res.Error != nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	err = bcrypt.CompareHashAndPassword([]byte(currentUser.PasswordDigest), []byte(params.Password))
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	sessionId, err := middlewares.StartSession(currentUser.ID)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	cookie, err := middlewares.GenerateJwt(sessionId)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	c.Response().Header().Set(echo.HeaderSetCookie, cookie)
	resp.ID = currentUser.ID
	resp.Name = currentUser.Name
	resp.Email = currentUser.Email
	return c.JSON(http.StatusOK, resp)
}
