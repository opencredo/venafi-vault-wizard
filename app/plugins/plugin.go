package plugins

import (
	"fmt"

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

// GetCatalogName returns the name of the plugin as it appears in the plugin catalog. This does not include the plugin
// version, to allow the plugin to be updated without needed to remount the associated instances. However it does
// include the mount path to allow the version of the plugin to vary independently between different mounted instances
// of it.
func (p *Plugin) GetCatalogName() string {
	return fmt.Sprintf("%s-%s", p.Type, p.MountPath)
}

// GetFileName returns the filename of the plugin as it will be found in the plugin directory. This includes the plugin
// version to allow different versions of the plugin to be present on the Vault server, and for the catalog entries to
// reference different ones depending on their mounts' use case.
func (p *Plugin) GetFileName() string {
	return fmt.Sprintf("%s_%s", p.Type, p.Version)
}
