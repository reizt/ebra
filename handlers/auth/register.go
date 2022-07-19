package auth

import (
	"net/http"

	"github.com/reizt/ebra/bindings"
	"github.com/reizt/ebra/models"
	"github.com/reizt/ebra/renderings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
)

var (
	bcryptCost = 12
)

func Register(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)
	params := &bindings.CreateUserRequest{}
	user := &models.User{}
	if err := (&echo.DefaultBinder{}).BindBody(c, &params); err != nil {
		return err
	}
	// Validation for now
	if params.Name == "" || params.Email == "" || params.Password == "" {
		return c.JSON(http.StatusBadRequest, renderings.NotFoundResponse{
			Message: "can't be blank",
		})
	}
	user.Name = params.Name
	user.Email = params.Email
	digest, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil {
		return err
	}
	user.PasswordDigest = string(digest)
	if err := db.Select("ID", "Name", "Email", "PasswordDigest").Create(&user).Error; err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, user)
}
