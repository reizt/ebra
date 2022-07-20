package users

import (
	"net/http"

	"github.com/reizt/ebra/conf"
	"github.com/reizt/ebra/models"
	"github.com/reizt/ebra/renderings"
	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
)

func GetUsers(c echo.Context) error {
	db := c.Get(conf.DbContextKey).(*gorm.DB)
	users := []renderings.User{} // Response will be null when initialize var users by new(), expects []
	db.Model(&models.User{}).Order("created_at desc").Limit(50).Find(&users)
	return c.JSON(http.StatusOK, users)
}
