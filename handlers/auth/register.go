package auth

import (
	"net/http"

	"github.com/reizt/ebra/bindings"
	"github.com/reizt/ebra/conf"
	"github.com/reizt/ebra/models"
	"github.com/reizt/ebra/renderings"

	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
)

func Register(c echo.Context) error {
	db := c.Get(conf.DbContextKey).(*gorm.DB)
	params := &bindings.CreateUserRequest{}
	resp := &renderings.UserResponse{}
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
	if err := db.Create(&user).Error; err != nil {
		return err
	}
	resp.ID = user.ID
	resp.Name = user.Name
	resp.Email = user.Email
	return c.JSON(http.StatusCreated, resp)
}
