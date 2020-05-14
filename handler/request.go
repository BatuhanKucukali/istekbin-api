package handler

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Request struct {
	Method      string            `json:"method"`
	Host        string            `json:"host"`
	UserAgent   string            `json:"user_agent"`
	Ip          string            `json:"ip"`
	Uri         string            `json:"uri"`
	ContentType string            `json:"content_type"`
	Header      map[string]string `json:"header"`
	Body        string            `json:"body"`
	CreatedAt   time.Time         `json:"created_at"`
}

const ForbiddenContentType = "multipart/form-data"

func RequestHandler(c echo.Context) error {
	uString := c.Param("uuid")
	u, err := uuid.Parse(uString)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid path.")
	}
	fmt.Println(u)

	req := c.Request()
	contentType := req.Header.Get("Content-Type")

	if isSupportedContentType(contentType) { // TODO it must be support all content type
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("%s not supported yet.", ForbiddenContentType))
	}

	r := new(Request)
	r.Method = req.Method
	r.Host = req.Host
	r.Uri = req.RequestURI
	r.ContentType = contentType
	r.UserAgent = req.UserAgent()
	r.Ip = req.RemoteAddr
	r.CreatedAt = time.Now()
	r.Header = getHeader(req.Header)

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "body can not reading.")
	}
	r.Body = string(body)

	return c.JSON(http.StatusOK, r)
}

func getHeader(header http.Header) map[string]string {
	var headerMap = make(map[string]string)
	for key, value := range header {
		headerMap[key] = strings.Join(value, ",")
	}
	return headerMap
}

func isSupportedContentType(contentType string) bool {
	return strings.Contains(contentType, ForbiddenContentType)
}
