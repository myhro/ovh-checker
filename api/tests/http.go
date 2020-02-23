package tests

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
)

// AuthHeader returns a base64-encoded Authorization token header
func AuthHeader(email, token string) string {
	pair := fmt.Sprintf("%v:%v", email, token)
	encoded := base64.StdEncoding.EncodeToString([]byte(pair))
	header := "Token " + encoded
	return header
}

// Get does an HTTP GET request against an HTTP handler
func Get(r http.Handler, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", path, nil)
	return request(r, req)
}

// GetWithHeaders does an HTTP GET request with custom headers against an HTTP handler
func GetWithHeaders(r http.Handler, path string, headers map[string]string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", path, nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return request(r, req)
}

// Post does an HTTP POST request against an HTTP handler
func Post(r http.Handler, path, body string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("POST", path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	return request(r, req)
}

func request(r http.Handler, req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
