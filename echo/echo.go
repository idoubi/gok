package echo

import (
	"github.com/idoubi/gok/echo/middleware"
	"github.com/labstack/echo/v4"
	em "github.com/labstack/echo/v4/middleware"
)

// New: new echo instance with common use middlewares
func New() *echo.Echo {
	e := echo.New()
	e.Validator = middleware.NewValidator()
	e.Use(em.Logger())
	e.Use(middleware.ApiContextWithConfig())
	e.Use(middleware.ErrorHandler)

	return e
}
