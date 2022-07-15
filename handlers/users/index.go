package users

import (
	"net/http"

	"github.com/reizt/ebra/models"
	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
)

func GetUsers(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)
	users := []models.User{} // Response will be null when initialize var users by new(), expects []
	db.Order("created_at desc").Limit(50).Find(&users)
	return c.JSON(http.StatusOK, users)
}
