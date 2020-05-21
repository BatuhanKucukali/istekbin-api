package handler

import (
	"github.com/labstack/echo"
	"net/http"
)

func HomeHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome :)")
}
