package mot

import (
	"github.com/idoubi/mot/middleware"
	"github.com/labstack/echo/v4"
	em "github.com/labstack/echo/v4/middleware"
)

// NewEcho: new echo instance
func NewEcho() *echo.Echo {
	e := echo.New()
	e.Validator = middleware.NewValidator()
	e.Use(em.Logger())
	e.Use(middleware.ApiContextWithConfig())
	e.Use(middleware.ErrorHandler)

	return e
}
