package main

import (
	"net/http"

	"github.com/reizt/ebra/config"
	"github.com/reizt/ebra/handlers/users"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	config.Migrate()
	db := config.ConnectMySQL()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(config.DbContextKey, db)
			return next(c)
		}
	})

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/users", users.GetUsers)
	e.POST("/users", users.CreateUser)
	e.GET("/users/:id", users.GetUserById)
	e.PATCH("/users/:id", users.UpdateUser)
	e.DELETE("/users/:id", users.DeleteUser)

	e.Logger.Fatal(e.Start(":3000"))
}
