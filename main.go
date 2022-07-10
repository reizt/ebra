package main

import (
	"net/http"
	"ebra/config"
	"ebra/controllers"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	config.Migrate()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/users", controllers.GetAllUsers)
	e.POST("/users", controllers.CreateUser)
	e.GET("/users/:id", controllers.GetUserById)
	e.PATCH("/users/:id", controllers.UpdateUser)
	e.DELETE("/users/:id", controllers.DeleteUser)

	e.Logger.Fatal(e.Start(":3000"))
}
