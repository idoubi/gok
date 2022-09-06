package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// ErrorHandler: custom error handler
func ErrorHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		c := ctx.(*ApiContext)

		err := next(c)

		if err != nil {
			if httpErr, ok := err.(*echo.HTTPError); ok {
				if httpErr.Code == http.StatusUnauthorized || httpErr.Message == "missing or malformed jwt" {
					return c.RespStd(-2, "unauthorized", "", nil)
				}
			}

			return c.RespErrWithDetail("system error", err.Error())
		}

		return nil
	}
}
