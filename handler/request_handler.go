package handler

import (
	"encoding/json"
	"fmt"
	"github.com/batuhankucukali/istekbin/config"
	"github.com/go-redis/redis/v7"
	"github.com/google/uuid"
	"github.com/labstack/echo"
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

func RequestHandler(conf *config.App, rd *redis.Client) func(c echo.Context) error {
	return func(c echo.Context) error {
		u, err := uuid.Parse(c.Param("uuid"))
		if err != nil {
			return echo.ErrNotFound
		}

		req := c.Request()

		isForbiddenHeaderExist, forbiddenHeader := isForbiddenHeaderExist(req.Header, *conf)
		if isForbiddenHeaderExist {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("%s header is forbidden.", forbiddenHeader))
		}

		reqVal, err := rd.Get(u.String()).Result()
		if err != nil {
			return echo.ErrNotFound
		}

		contentType := req.Header.Get(echo.HeaderContentType)

		r := new(Request)
		r.Method = req.Method
		r.Host = req.Host
		r.Uri = getUri(req, u)
		r.ContentType = contentType
		r.UserAgent = req.UserAgent()
		r.Ip = getIp(conf, req)
		r.CreatedAt = time.Now()
		r.Headers = getHeaders(req.Header)
		r.Cookies = getCookies(req)

		if isMultipartForm(contentType) {
			body, err := parseMultipartBody(c)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "body can not reading.")
			}
			r.Body = body
			r.ContentType = echo.MIMEMultipartForm
		} else {
			body, err := parseBody(req.Body)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "body can not reading.")
			}
			r.Body = body
		}
		var rl []Request

		if len(reqVal) > 0 {
			if err := json.Unmarshal([]byte(reqVal), &rl); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "request can not deserialized.")
			}
		}

		requests := append([]Request{*r}, rl...)

		if err := set(conf, rd, u.String(), requests); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "request can not serialized.")
		}

		return c.String(http.StatusOK, "ok")
	}
}

func isForbiddenHeaderExist(headers http.Header, conf config.App) (bool, string) {
	for key := range headers {
		for _, forbiddenHeader := range conf.ForbiddenHeaders {
			if key == forbiddenHeader {
				return true, forbiddenHeader
			}
		}
	}
	return false, ""
}

func getCookies(req *http.Request) map[string]string {
	var cookieMap = make(map[string]string)
	for _, c := range req.Cookies() {
		cookieMap[c.Name] = c.Value
	}
	return cookieMap
}

func getIp(conf *config.App, req *http.Request) string {
	if conf.Env == "production" {
		return req.Header.Get("X-Forwarded-For")
	}
	return req.RemoteAddr
}

func getUri(req *http.Request, u uuid.UUID) string {
	return strings.Replace(req.RequestURI, "r/"+u.String()+"/", "", -1)
}

func set(conf *config.App, rd *redis.Client, key string, value interface{}) error {
	p, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return rd.Set(key, p, conf.RequestStoreTime).Err()
}

func getHeaders(header http.Header) map[string]string {
	var headerMap = make(map[string]string)
	for key, value := range header {
		headerMap[key] = strings.Join(value, ",")
	}
	return headerMap
}

func parseBody(reqBody io.ReadCloser) (string, error) {
	body, err := ioutil.ReadAll(reqBody)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func parseMultipartBody(c echo.Context) (string, error) {
	multi, err := c.MultipartForm()

	if err != nil {
		return "", err
	}

	var body string

	for form := range multi.Value {
		body += fmt.Sprintf("%s=%s\n", form, c.FormValue(form))
	}

	return body, nil
}

func isMultipartForm(contentType string) bool {
	return strings.Contains(contentType, echo.MIMEMultipartForm)
}
