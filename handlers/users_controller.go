package handlers

import (
	"net/http"

	"github.com/reizt/ebra/models"
	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
)

func GetAllUsers(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)
	users := []models.User{} // Response will be null when initialize var users by new(), expects []
	db.Order("created_at desc").Limit(50).Find(&users)
	return c.JSON(http.StatusOK, users)
}
func GetUserById(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)
	user := new(models.User)
	id := c.Param("id")
	if err := db.First(&user, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"message": "user not found",
		})
	}
	return c.JSON(http.StatusOK, user)
}
func CreateUser(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)
	user := &models.User{}
	if err := (&echo.DefaultBinder{}).BindBody(c, &user); err != nil {
		return err
	}
	// Validation for now
	if user.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "name can't be blank",
		})
	}
	if err := db.Select("ID", "Name").Create(&user).Error; err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, user)
}
func UpdateUserById(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)
	user := new(models.User)
	id := c.Param("id")
	findRes := db.First(&user, "id = ?", id)

	if findRes.Error != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"message": "user not found",
		})
	}

	err := c.Bind(&user)
	if err != nil {
		return err
	}
	updateRes := db.Model(&user).Updates(models.User{
		Name: user.Name,
	})
	if updateRes.Error != nil {
		return updateRes.Error
	}
	return c.NoContent(http.StatusNoContent)
}
func DeleteUserById(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)
	user := new(models.User)
	findRes := db.First(&user, "id = ?", c.Param("id"))

	if findRes.Error != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"message": "user not found",
		})
	}
	deleteRes := db.Delete(&user)
	if deleteRes.Error != nil {
		return deleteRes.Error
	}
	return c.NoContent(http.StatusNoContent)
}
