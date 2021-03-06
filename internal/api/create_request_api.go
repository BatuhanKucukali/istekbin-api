package api

import (
	"encoding/json"
	"fmt"
	"github.com/batuhankucukali/istekbin-api/internal/config"
	"github.com/go-redis/redis/v7"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"io"
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
	Headers     map[string]string `json:"headers"`
	Cookies     map[string]string `json:"cookies"`
	Body        string            `json:"body"`
	CreatedAt   time.Time         `json:"created_at"`
}

// CreateRequest
// @Summary Create request - EXAMPLE!!! - Swagger does not allowed multiple http method.
// @Description Route accept all of http methods. Swagger does not allowed multiple http method.
// @Success 200 {string} string "ok"
// @Param uuid path string true "uuid"
// @Param body body string false "body"
// @Param hello formData string false "world"
// @Param hello header string false "world"
// @Param hello query string false "world"
// @Failure 400 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /r/{uuid} [post]
func CreateRequest(conf *config.App, rd *redis.Client) func(c echo.Context) error {
	return func(c echo.Context) error {
		u, err := uuid.Parse(c.Param("uuid"))
		if err != nil {
			return echo.ErrNotFound
		}

		result, err := rd.Get(u.String()).Result()
		if err != nil {
			return echo.ErrNotFound
		}

		var rl []Request
		if len(result) > 0 {
			if err := json.Unmarshal([]byte(result), &rl); err != nil {
				log.Errorf("deserialize error. %s", err)
				return echo.NewHTTPError(http.StatusInternalServerError, "deserialize error")
			}
		}

		if len(rl) >= conf.MaxRequestCount {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("max request count is %d", conf.MaxRequestCount))
		}

		req := c.Request()

		contentType := req.Header.Get(echo.HeaderContentType)

		r := new(Request)
		r.Method = req.Method
		r.Host = req.Host
		r.Uri = getUri(req, u)
		r.ContentType = contentType
		r.UserAgent = req.UserAgent()
		r.Ip = c.RealIP()
		r.CreatedAt = time.Now()
		r.Headers = getHeaders(req.Header, *conf)
		r.Cookies = getCookies(req)

		body, err := getBody(req.Body)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "body can not reading.")
		}
		r.Body = body

		requests := append([]Request{*r}, rl...)

		if err := set(conf, rd, u.String(), requests); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "request can not serialized.")
		}

		return c.String(http.StatusOK, "ok")
	}
}

func isForbiddenHeader(header string, conf config.App) bool {
	if len(conf.ForbiddenHeaders) == 0 {
		return false
	}

	for _, forbiddenHeader := range conf.ForbiddenHeaders {
		if header == forbiddenHeader {
			return true
		}
	}

	return false
}

func getCookies(req *http.Request) map[string]string {
	var cookieMap = make(map[string]string)
	for _, c := range req.Cookies() {
		cookieMap[c.Name] = c.Value
	}
	return cookieMap
}

func getUri(req *http.Request, u uuid.UUID) string {
	return strings.Replace(req.RequestURI, "/r/"+u.String(), "", 1)
}

func set(conf *config.App, rd *redis.Client, key string, value interface{}) error {
	requestBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return rd.Set(key, requestBytes, conf.RequestStoreTime).Err()
}

func getHeaders(header http.Header, conf config.App) map[string]string {
	var headerMap = make(map[string]string)
	for key, value := range header {
		if !isForbiddenHeader(key, conf) {
			headerMap[key] = strings.Join(value, ",")
		}
	}
	return headerMap
}

func getBody(reqBody io.ReadCloser) (string, error) {
	body, err := ioutil.ReadAll(reqBody)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
