package helpers

import (
	"io"
	"net/http"
	"net/http/httptest"

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
