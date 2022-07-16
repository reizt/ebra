package users_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/reizt/ebra/conf"
	"github.com/reizt/ebra/models"

	"github.com/labstack/echo/v4"
)

func InitTestContext(method string, path string, body io.Reader) (ctx echo.Context, req *http.Request, rec *httptest.ResponseRecorder) {
	req = httptest.NewRequest(http.MethodGet, "/", body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	ctx = echo.New().NewContext(req, rec)
	ctx.SetPath(path)
	return ctx, req, rec
}

func TestMain(m *testing.M) {
	db := conf.ConnectMySQL()
	db.AutoMigrate(&models.User{})
	m.Run()
}
