package main

import (
	"net/http"

	"github.com/reizt/ebra/handlers"

	"github.com/reizt/ebra/conf"
	"github.com/reizt/ebra/handlers/users"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	conf.LoadEnv()
	e := echo.New()
	conf.Migrate()
	db := conf.ConnectMySQL()
	e.Use(middleware.LoggerWithConfig(
		middleware.LoggerConfig{
			Format: "\x1b[32m${method}\x1b[0m \x1b[32m${status}\x1b[0m \x1b[32m${uri}\x1b[0m\n",
		},
	))
	e.Static("/", "public")
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(conf.DbContextKey, db)
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
	g := e.Group("/test")
	g.POST("/mail", handlers.SendTestMail)

	e.Logger.Fatal(e.Start(":3000"))
}
