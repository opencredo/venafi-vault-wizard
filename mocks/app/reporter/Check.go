// Code generated by mockery v2.10.4. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Check is an autogenerated mock type for the Check type
type Check struct {
	mock.Mock
}

// Error provides a mock function with given fields: message
func (_m *Check) Error(message string) {
	_m.Called(message)
}

// Errorf provides a mock function with given fields: status, a
func (_m *Check) Errorf(status string, a ...interface{}) {
	_m.Called(status, a)
}

// Success provides a mock function with given fields: message
func (_m *Check) Success(message string) {
	_m.Called(message)
}

// Successf provides a mock function with given fields: status, a
func (_m *Check) Successf(status string, a ...interface{}) {
	_m.Called(status, a)
}

// UpdateStatus provides a mock function with given fields: status
func (_m *Check) UpdateStatus(status string) {
	_m.Called(status)
}

// UpdateStatusf provides a mock function with given fields: status, a
func (_m *Check) UpdateStatusf(status string, a ...interface{}) {
	_m.Called(status, a)
}

// Warning provides a mock function with given fields: message
func (_m *Check) Warning(message string) {
	_m.Called(message)
}

// Warningf provides a mock function with given fields: status, a
func (_m *Check) Warningf(status string, a ...interface{}) {
	_m.Called(status, a)
}
