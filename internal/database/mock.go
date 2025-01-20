package database

import (
	"database/sql"

	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Connect() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockDB) CheckHealth() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	call := m.Called(query, args)
	return call.Get(0).(*sql.Rows), call.Error(1)
}

func (m *MockDB) QueryRow(query string, args ...interface{}) *sql.Row {
	argsMock := m.Called(query, args)
	return argsMock.Get(0).(*sql.Row)
}

func (m *MockDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	call := m.Called(query, args)
	return call.Get(0).(sql.Result), call.Error(1)
}
