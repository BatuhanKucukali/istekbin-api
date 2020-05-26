package api

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// Home
// @Summary Welcome page
// @Accept  plain
// @Success 200 {string} string ""
// @Router / [get]
func Home(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome :)")
}
