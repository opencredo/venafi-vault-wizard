package plugins

import (
	"github.com/hashicorp/hcl/v2"

	"github.com/opencredo/venafi-vault-wizard/app/reporter"
	"github.com/opencredo/venafi-vault-wizard/app/vault/api"
)

// Generic wrapper around a specific plugin implementation representing concerns common to all Vault plugins
type Plugin struct {
	// Type, the first label should specify which plugin the block refers to
	Type string `hcl:"type,label"`
	// MountPath, the second label should specify the Vault mount path
	MountPath string `hcl:"mount_path,label"`
	// Version is the version of the plugin to use, specified as the Git tag of the release
	Version string `hcl:"version"`
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
