package middleware

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// ApiContext: custom context
type ApiContext struct {
	echo.Context
}

// resp: custom response body
type resp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Valid: valid request params
func (c *ApiContext) Valid(req interface{}) error {
	if err := c.Bind(req); err != nil {
		if v, ok := err.(*echo.HTTPError); ok {
			return fmt.Errorf("%s", v.Message)
		}

		return err
	}
	if err := c.Validate(req); err != nil {
		return err
	}

	return nil
}

// RespOk: success response
func (c *ApiContext) RespOk(msg string) error {
	return c.RespStd(0, msg, nil)
}

// RespOkWithData: success response with data
func (c *ApiContext) RespOkWithData(msg string, data interface{}) error {
	return c.RespStd(0, msg, data)
}

// RespErr: fail response
func (c *ApiContext) RespErr(msg string) error {
	return c.RespStd(-1, msg, nil)
}

// RespStd: standard response
func (c *ApiContext) RespStd(code int, msg string, data interface{}) error {
	return c.JSON(http.StatusOK, resp{code, msg, data})
}

// ApiContextWithConfig: custom context middleware
func ApiContextWithConfig() echo.MiddlewareFunc {
	return apiContext
}

func apiContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		c := &ApiContext{ctx}

		return next(c)
	}
}