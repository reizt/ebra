package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/reizt/ebra/middlewares"
)

func Signout(c echo.Context) error {
	token := c.Request().Header[echo.HeaderAuthorization]
	// Returns 400 if authorization header doesn't exists
	if len(token) == 0 {
		return echo.ErrBadRequest
	}

	claims, err := middlewares.ValidateJwt(token[0])
	// Returns 401 if requested JWT is invalid
	if err != nil {
		return echo.ErrUnauthorized
	}

	sessionId := claims["sessionId"].(string)
	if err := middlewares.DisableSession(sessionId); err != nil {
		return echo.ErrBadRequest
	}
	return c.NoContent(http.StatusNoContent)
}
