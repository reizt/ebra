package users

import (
	"net/http"

	"github.com/reizt/ebra/conf"
	"github.com/reizt/ebra/models"
	"github.com/reizt/ebra/renderings"
	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
)

func GetUserById(c echo.Context) error {
	db := c.Get(conf.DbContextKey).(*gorm.DB)
	user := new(models.User)
	id := c.Param("id")
	if err := db.First(&user, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, renderings.NotFoundResponse{
			Message: "user not found",
		})
	}
	return c.JSON(http.StatusOK, user)
}
