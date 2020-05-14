package handler

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"net/http"
)

func CreateHandler(c echo.Context) error {
	u := uuid.New()
	rPath := fmt.Sprintf("r/%s", u.String()) // TODO add baseUrl

	c.Response().Header().Add("Location", rPath)
	return c.NoContent(http.StatusCreated)
}
