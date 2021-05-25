// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	hcl "github.com/hashicorp/hcl/v2"
	api "github.com/opencredo/venafi-vault-wizard/app/vault/api"

	hclwrite "github.com/hashicorp/hcl/v2/hclwrite"

	mock "github.com/stretchr/testify/mock"

	plugins "github.com/opencredo/venafi-vault-wizard/app/plugins"

	reporter "github.com/opencredo/venafi-vault-wizard/app/reporter"
)

// PluginImpl is an autogenerated mock type for the PluginImpl type
type PluginImpl struct {
	mock.Mock
}

// Check provides a mock function with given fields: report, vaultClient
func (_m *PluginImpl) Check(report reporter.Report, vaultClient api.VaultAPIClient) error {
	ret := _m.Called(report, vaultClient)

	var r0 error
	if rf, ok := ret.Get(0).(func(reporter.Report, api.VaultAPIClient) error); ok {
		r0 = rf(report, vaultClient)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Configure provides a mock function with given fields: report, vaultClient
func (_m *PluginImpl) Configure(report reporter.Report, vaultClient api.VaultAPIClient) error {
	ret := _m.Called(report, vaultClient)

	var r0 error
	if rf, ok := ret.Get(0).(func(reporter.Report, api.VaultAPIClient) error); ok {
		r0 = rf(report, vaultClient)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GenerateConfigAndWriteHCL provides a mock function with given fields: hclBody
func (_m *PluginImpl) GenerateConfigAndWriteHCL(hclBody *hclwrite.Body) error {
	ret := _m.Called(hclBody)

	var r0 error
	if rf, ok := ret.Get(0).(func(*hclwrite.Body) error); ok {
		r0 = rf(hclBody)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetDownloadURL provides a mock function with given fields:
func (_m *PluginImpl) GetDownloadURL() (string, error) {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ParseConfig provides a mock function with given fields: config, evalContext
func (_m *PluginImpl) ParseConfig(config *plugins.Plugin, evalContext *hcl.EvalContext) error {
	ret := _m.Called(config, evalContext)

	var r0 error
	if rf, ok := ret.Get(0).(func(*plugins.Plugin, *hcl.EvalContext) error); ok {
		r0 = rf(config, evalContext)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ValidateConfig provides a mock function with given fields:
func (_m *PluginImpl) ValidateConfig() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
