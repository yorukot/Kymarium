package testutil

import (
	"io"
	"net/http/httptest"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/yorukot/kymarium/api/middleware"
)

// NewEchoContext creates a fresh Echo context and recorder.
func NewEchoContext(method, path string, body io.Reader) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(method, path, body)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}

// Authenticate sets the user ID on the context to mimic middleware.
func Authenticate(c echo.Context, userID int64) {
	c.Set(string(middleware.UserIDKey), strconv.FormatInt(userID, 10))
}

// SetJSONHeader ensures the request content-type is JSON.
func SetJSONHeader(c echo.Context) {
	c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
}
