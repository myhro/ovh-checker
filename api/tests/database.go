package tests

import (
	"database/sql"

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
