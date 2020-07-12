package controllers

import (
	"net/http"

	"github.com/labstack/echo"
)

// Hello handler method just returns hello
func Hello(c echo.Context) (err error) {

	return c.String(http.StatusOK, "hello")
}
