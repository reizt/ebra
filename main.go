package main

import (
	"net/http"
	"werp/api/config"
	"werp/api/controllers"

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

	e.Logger.Fatal(e.Start(":3000"))
}
