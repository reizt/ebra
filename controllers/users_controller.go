package controllers

import (
	"net/http"
	"werp/api/config"
	"werp/api/models"

	"github.com/labstack/echo/v4"
)

var db = config.ConnectMySQL()

func GetAllUsers(c echo.Context) error {
	var users []models.User
	db.Limit(50).Find(&users)
	return c.JSON(http.StatusOK, users)
}
func GetUserById(c echo.Context) error {
	user := new(models.User)
	id := c.Param("id")
	db.Limit(1).Find(&user, "id = ?", id)
	if user != new(models.User) {
		return c.JSON(http.StatusOK, user)
	} else {
		return c.JSON(http.StatusNotFound, map[string]string{
			"message": "User not found",
		})
	}
}
func CreateUser(c echo.Context) error {
	user := new(models.User)
	if err := c.Bind(&user); err != nil {
		return err
	}

	result := db.Select("ID", "Name").Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return c.JSON(http.StatusCreated, user)
}
func UpdateUser(c echo.Context) error {
	user := new(models.User)
	result := db.First(&user, "id = ?", c.Param("id"))

	if result.Error != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"message": "User not found",
		})
	}

	err := c.Bind(&user)
	if err != nil {
		return err
	}
	res := db.Model(&user).Updates(models.User{
		Name: user.Name,
	})
	if res.Error != nil {
		return res.Error
	}
	return c.JSON(http.StatusAccepted, user)
}
