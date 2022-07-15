package main

import (
	"net/http"

	"github.com/reizt/ebra/config"
	"github.com/reizt/ebra/handlers"

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
	e.GET("/users", handlers.GetAllUsers)
	e.POST("/users", handlers.CreateUser)
	e.GET("/users/:id", handlers.GetUserById)
	e.PATCH("/users/:id", handlers.UpdateUserById)
	e.DELETE("/users/:id", handlers.DeleteUserById)

	e.Logger.Fatal(e.Start(":3000"))
}
