// Code generated by mockery v2.10.4. DO NOT EDIT.

package mocks

import (
	io "io"

	mock "github.com/stretchr/testify/mock"
)

// VaultSSHClient is an autogenerated mock type for the VaultSSHClient type
type VaultSSHClient struct {
	mock.Mock
}

// AddIPCLockCapabilityToFile provides a mock function with given fields: filename
func (_m *VaultSSHClient) AddIPCLockCapabilityToFile(filename string) error {
	ret := _m.Called(filename)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(filename)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CheckOSArch provides a mock function with given fields:
func (_m *VaultSSHClient) CheckOSArch() (string, string, error) {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 string
	if rf, ok := ret.Get(1).(func() string); ok {
		r1 = rf()
	} else {
		r1 = ret.Get(1).(string)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func() error); ok {
		r2 = rf()
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Close provides a mock function with given fields:
func (_m *VaultSSHClient) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FileExists provides a mock function with given fields: filepath
func (_m *VaultSSHClient) FileExists(filepath string) (bool, error) {
	ret := _m.Called(filepath)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(filepath)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(filepath)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsIPCLockCapabilityOnFile provides a mock function with given fields: filename
func (_m *VaultSSHClient) IsIPCLockCapabilityOnFile(filename string) (bool, error) {
	ret := _m.Called(filename)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(filename)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(filename)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WriteFile provides a mock function with given fields: sourceFile, hostDestination
func (_m *VaultSSHClient) WriteFile(sourceFile io.Reader, hostDestination string) error {
	ret := _m.Called(sourceFile, hostDestination)

	var r0 error
	if rf, ok := ret.Get(0).(func(io.Reader, string) error); ok {
		r0 = rf(sourceFile, hostDestination)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
