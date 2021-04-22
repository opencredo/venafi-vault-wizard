package plugins

import (
	"errors"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"

	pki_backend "github.com/opencredo/venafi-vault-wizard/app/plugins/venafi/pki-backend"
	"github.com/opencredo/venafi-vault-wizard/app/reporter"
	"github.com/opencredo/venafi-vault-wizard/app/vault/api"
)

var ErrPluginNotFound = errors.New("plugin not found")

// Generic wrapper around a specific plugin implementation representing concerns common to all Vault plugins
type Plugin struct {
	// Type, the first label should specify which plugin the block refers to
	Type string `hcl:"type,label"`
	// MountPath, the second label should specify the Vault mount path
	MountPath string `hcl:"mount_path,label"`
	// Version is an optional attribute that should default to the latest version if not specified
	Version string `hcl:"version,optional"`
	// Config allows the rest of the plugin body to be decoded separately
	Config hcl.Body `hcl:",remain"`

	// Impl is an implementation of the Plugin interface, defining both Configure and Check methods to perform the
	// relevant Vault configuration tasks for the specific plugin. It is not populated by the initial HCL decoding, as
	// the schema unique to each plugin, so it is instead populated after the fact by NewConfig, which calls
	// plugins.LookupPlugin to find an implementation based on the Type field.
	Impl PluginImpl
}

// PluginImpl is what each plugin struct should implement, along with having its config schema defined with its struct
// fields and the relevant HCL struct tags.
type PluginImpl interface {
	// GetDownloadURL returns a URL to download the required version of the plugin
	GetDownloadURL() (string, error)
	// Configure makes the necessary changes to Vault to configure the plugin
	Configure(report reporter.Report, vaultClient api.VaultAPIClient) error
	// Check is similar to Configure, except it shouldn't make any changes, only validate what is already there
	Check(report reporter.Report, vaultClient api.VaultAPIClient) error
	// ValidateConfig performs validation of the supplied configuration data, specific to the plugin
	ValidateConfig() error
}

// Map of "constructors" which maps all the supported plugin types to their associated PluginImpl implementations.
// Uses a function to force each instance to be a copy, and to allow injection of fields from the wrapper Plugin struct
type pluginImplConstructor func(*Plugin) PluginImpl

var supportedPlugins = map[string]pluginImplConstructor{
	"venafi-pki-backend": func(config *Plugin) PluginImpl {
		return &pki_backend.VenafiPKIBackendConfig{
			MountPath: config.MountPath,
			Version:   config.Version,
		}
	},
}

// LookupPlugin goes from a generic config.Plugin and looks up its specific PluginImpl based on the Type
// field. It takes a pointer to an hcl.EvalContext in order to provide the same config functions available to the rest
// of the config, and in future maybe some global variables, and then decodes the plugin-specific part of the plugin
// block.
func LookupPlugin(config *Plugin, evalContext *hcl.EvalContext) (PluginImpl, error) {
	constructor, ok := supportedPlugins[config.Type]
	if !ok {
		return nil, ErrPluginNotFound
	}

	plugin := constructor(config)

	diagnostics := gohcl.DecodeBody(config.Config, evalContext, plugin)
	if diagnostics.HasErrors() {
		return nil, diagnostics
	}

	return plugin, nil
}
