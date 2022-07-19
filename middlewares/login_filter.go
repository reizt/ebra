package middlewares

import (
	"fmt"
	"regexp"

	"github.com/labstack/echo/v4"
)

func SigninFilter(next echo.HandlerFunc) echo.HandlerFunc {
	skipper := func(c echo.Context) bool {
		path := c.Request().URL.Path
		matched, err := regexp.Match("/auth/(signin|register)", []byte(path))
		if err != nil {
			fmt.Println(err)
		}
		return matched && err == nil
	}
	return func(c echo.Context) error {
		// Skip if request doesn't need to verify login
		if skipper(c) {
			return next(c)
		}
		token := c.Request().Header[echo.HeaderAuthorization]
		// Returns 400 if authorization header doesn't exists
		if len(token) == 0 {
			return echo.ErrBadRequest
		}

		claims, err := ValidateJwt(token[0])
		// Returns 401 if requested JWT is invalid
		if err != nil {
			return echo.ErrUnauthorized
		}

		sessionId := claims["sessionId"].(string)
		currentUser, err := GetCurrentUser(sessionId)
		// Returns 401 if current user could not be found
		if err != nil || currentUser == nil {
			return echo.ErrUnauthorized
		}
		return next(c)
	}
}
