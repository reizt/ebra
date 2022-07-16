package users

import (
	"net/http"

	"github.com/reizt/ebra/bindings"
	"github.com/reizt/ebra/models"
	"github.com/reizt/ebra/renderings"

	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
)

func CreateUser(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)
	params := &bindings.CreateUserRequest{}
	user := &models.User{}
	if err := (&echo.DefaultBinder{}).BindBody(c, &params); err != nil {
		return err
	}
	user.Name = params.Name
	// Validation for now
	if user.Name == "" {
		return c.JSON(http.StatusBadRequest, renderings.NotFoundResponse{
			Message: "name can't be blank",
		})
	}
	if err := db.Select("ID", "Name").Create(&user).Error; err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, user)
}
