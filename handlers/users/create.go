package users

import (
	"net/http"

	"github.com/reizt/ebra/bindings"
	"github.com/reizt/ebra/conf"
	"github.com/reizt/ebra/models"
	"github.com/reizt/ebra/renderings"

	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
)

func CreateUser(c echo.Context) error {
	db := c.Get(conf.DbContextKey).(*gorm.DB)
	params := &bindings.CreateUserRequest{}
	resp := &renderings.User{}
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
	resp.ID = user.ID
	resp.Name = user.Name
	resp.Email = user.Email
	resp.CreatedAt = user.CreatedAt
	resp.UpdatedAt = user.UpdatedAt
	return c.JSON(http.StatusCreated, resp)
}
