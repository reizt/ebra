package users

import (
	"net/http"

	"github.com/reizt/ebra/renderings"

	"github.com/reizt/ebra/models"
	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
)

func CreateUser(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)
	user := &models.User{}
	if err := (&echo.DefaultBinder{}).BindBody(c, &user); err != nil {
		return err
	}
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
