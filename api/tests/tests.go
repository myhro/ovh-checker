package tests

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/stretchr/testify/mock"
)

// MockedDatabase is a dummy database handler
type MockedDatabase struct {
	mock.Mock
}

// Exec mocks the same sql.DB method
func (m *MockedDatabase) Exec(q string, a ...interface{}) (sql.Result, error) {
	var res sql.Result
	args := m.Called(q, a)
	return res, args.Error(1)
}

// Get mocks the same sqlx.DB method
func (m *MockedDatabase) Get(d interface{}, q string, a ...interface{}) error {
	args := m.Called(d, q, a)
	return args.Error(0)
}

// Select mocks the same sqlx.DB method
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
