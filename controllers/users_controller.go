package controllers

import (
	"net/http"
	"echo-basic-rest-api/config"
	"echo-basic-rest-api/models"

	"github.com/labstack/echo/v4"
)

var db = config.ConnectMySQL()

func GetAllUsers(c echo.Context) error {
	users := new([]models.User)
	db.Limit(50).Find(&users)
	return c.JSON(http.StatusOK, users)
}
func GetUserById(c echo.Context) error {
	user := new(models.User)
	id := c.Param("id")
	getRes := db.First(&user, "id = ?", id)
	if getRes.Error != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"message": "User not found",
		})
	}
	return c.JSON(http.StatusOK, user)
}
func CreateUser(c echo.Context) error {
	user := new(models.User)
	if err := c.Bind(&user); err != nil {
		return err
	}

	createRes := db.Select("ID", "Name").Create(&user)
	if createRes.Error != nil {
		return createRes.Error
	}
	return c.JSON(http.StatusCreated, user)
}
func UpdateUser(c echo.Context) error {
	user := new(models.User)
	id := c.Param("id")
	findRes := db.First(&user, "id = ?", id)

	if findRes.Error != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"message": "User not found",
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
	return c.NoContent(http.StatusOK)
}
func DeleteUser(c echo.Context) error {
	user := new(models.User)
	findRes := db.First(&user, "id = ?", c.Param("id"))

	if findRes.Error != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"message": "User not found",
		})
	}
	deleteRes := db.Delete(&user)
	if deleteRes.Error != nil {
		return deleteRes.Error
	}
	return c.NoContent(http.StatusOK)
}
