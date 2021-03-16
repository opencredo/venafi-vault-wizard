// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	reporter "github.com/opencredo/venafi-vault-wizard/app/reporter"
	mock "github.com/stretchr/testify/mock"
)

// Report is an autogenerated mock type for the Report type
type Report struct {
	mock.Mock
}

// AddSection provides a mock function with given fields: name
func (_m *Report) AddSection(name string) reporter.Section {
	ret := _m.Called(name)

	var r0 reporter.Section
	if rf, ok := ret.Get(0).(func(string) reporter.Section); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(reporter.Section)
		}
	}

	return r0
}

// Finish provides a mock function with given fields: summary, message
func (_m *Report) Finish(summary string, message string) {
	_m.Called(summary, message)
}
