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

func UpdateUser(c echo.Context) error {
	db := c.Get(conf.DbContextKey).(*gorm.DB)
	user := &models.User{}
	id := c.Param("id")
	findRes := db.First(&user, "id = ?", id)

	if findRes.Error != nil {
		return c.JSON(http.StatusNotFound, renderings.NotFoundResponse{
			Message: "user not found",
		})
	}
	params := &bindings.UpdateUserRequest{}
	if err := (&echo.DefaultBinder{}).BindBody(c, &params); err != nil {
		return err
	}
	updateRes := db.Model(&user).Updates(models.User{
		Name: params.Name,
	})
	if updateRes.Error != nil {
		return updateRes.Error
	}
	return c.NoContent(http.StatusNoContent)
}
