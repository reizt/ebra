package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/reizt/ebra/conf"
)

func SetContexts() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		db := conf.ConnectMySQL()
		return func(c echo.Context) error {
			c.Set(conf.DbContextKey, db)
			return next(c)
		}
	}
}
