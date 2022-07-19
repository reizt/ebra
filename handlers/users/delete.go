package users

import (
	"net/http"

	"github.com/reizt/ebra/conf"
	"github.com/reizt/ebra/models"
	"github.com/reizt/ebra/renderings"
	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
)

func DeleteUser(c echo.Context) error {
	db := c.Get(conf.DbContextKey).(*gorm.DB)
	user := new(models.User)
	findRes := db.First(&user, "id = ?", c.Param("id"))

	if findRes.Error != nil {
		return c.JSON(http.StatusNotFound, renderings.NotFoundResponse{
			Message: "user not found",
		})
	}
	deleteRes := db.Delete(&user)
	if deleteRes.Error != nil {
		return deleteRes.Error
	}
	return c.NoContent(http.StatusNoContent)
}
