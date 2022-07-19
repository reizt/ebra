package main

import (
	"github.com/reizt/ebra/conf"
	"github.com/reizt/ebra/handlers"
	"github.com/reizt/ebra/handlers/auth"
	"github.com/reizt/ebra/handlers/users"
	"github.com/reizt/ebra/middlewares"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	conf.LoadEnv()
	conf.Migrate()
	e := echo.New()

	e.Use(middleware.LoggerWithConfig(
		middleware.LoggerConfig{
			Format: "${time_rfc3339_nano} ${host} \x1b[32m${method}\x1b[0m \x1b[32m${status}\x1b[0m \x1b[32m${uri}\x1b[0m\n",
		},
	))
	e.Use(middleware.Recover())
	e.Use(middlewares.SetContexts())
	e.Use(middlewares.SigninFilter())

	e.GET("/", handlers.Root)
	e.Static("/", "public")

	g := e.Group("/auth")
	g.POST("/signin", auth.Signin)
	g.DELETE("/signout", auth.Signout)
	g.POST("/register", auth.Register)

	e.GET("/users", users.GetUsers)
	e.POST("/users", users.CreateUser)
	e.GET("/users/:id", users.GetUserById)
	e.PATCH("/users/:id", users.UpdateUser)
	e.DELETE("/users/:id", users.DeleteUser)

	e.POST("/test/mail", handlers.SendTestMail)

	e.Logger.Fatal(e.Start(":3000"))
}
