package tests

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/stretchr/testify/mock"
)

// MockedDatabase is a dummy database handler
type MockedDatabase struct {
	mock.Mock
}

// Select mocks the sqlx.DB Select() method
func (m *MockedDatabase) Select(d interface{}, q string, a ...interface{}) error {
	args := m.Called(d, q, a)
	return args.Error(0)
}

// Get does an HTTP GET request against an HTTP handler
func Get(r http.Handler, path string) *httptest.ResponseRecorder {
	return request(r, "GET", path, nil)
}

// Post does an HTTP POST request against an HTTP handler
func Post(r http.Handler, path, body string) *httptest.ResponseRecorder {
	return request(r, "POST", path, strings.NewReader(body))
}

func request(r http.Handler, method, path string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
