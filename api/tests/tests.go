package tests

import (
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
	req, _ := http.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// Post does an HTTP POST request against an HTTP handler
func Post(r http.Handler, path, body string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("POST", path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
