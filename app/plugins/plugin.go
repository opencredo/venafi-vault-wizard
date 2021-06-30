package plugins

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/opencredo/venafi-vault-wizard/app/questions"
	"github.com/zclconf/go-cty/cty"

	"github.com/opencredo/venafi-vault-wizard/app/reporter"
	"github.com/opencredo/venafi-vault-wizard/app/vault/api"
)

// PluginConfig is a generic wrapper around a specific plugin implementation representing concerns common to all Vault plugins
type PluginConfig struct {
	// Type, the first label should specify which plugin the block refers to
	Type string `hcl:"type,label"`
	// MountPath, the second label should specify the Vault mount path
	MountPath string `hcl:"mount_path,label"`
	// Version is the version of the plugin to use, specified as the Git tag of the release
	Version string `hcl:"version"`
	// Filename is an optional field overriding the filename of the plugin binary in the plugin directory
	Filename string `hcl:"filename,optional"`
	// Config allows the rest of the plugin body to be decoded separately
	Config hcl.Body `hcl:",remain"`

	// VaultPluginType refers to whether the plugin is a secrets backend, or a database backend plugin
	VaultPluginType api.PluginType
	// Impl is an implementation of the Plugin interface, defining both Configure and Check methods to perform the
	// relevant Vault configuration tasks for the specific plugin. It is not populated by the initial HCL decoding, as
	// the schema unique to each plugin, so it is instead populated after the fact by NewConfig, which calls
	// plugins.LookupPlugin to find an implementation based on the Type field.
	Impl Plugin
}

// Plugin is what each plugin struct should implement, along with having its config schema defined with its struct
// fields and the relevant HCL struct tags.
type Plugin interface {
	// ParseConfig parses the PluginConfig.Config field (an hcl.Body) and populates its implementation-specific config struct.
	// It takes a pointer to an hcl.EvalContext in order to provide the same config functions available to the rest of
	// the config, and in future maybe some global variables, and then decodes the plugin-specific part of the plugin
	// block.
	ParseConfig(config *PluginConfig, evalContext *hcl.EvalContext) error
	// DownloadPlugin returns a byte slice with the plugin binary itself and a string with its SHA256
	DownloadPlugin() ([]byte, string, error)
	// Configure makes the necessary changes to Vault to configure the plugin
	Configure(report reporter.Report, vaultClient api.VaultAPIClient) error
	// Check is similar to Configure, except it shouldn't make any changes, only validate what is already there
	Check(report reporter.Report, vaultClient api.VaultAPIClient) error
	// ValidateConfig performs validation of the supplied configuration data, specific to the plugin
	ValidateConfig() error
	// GenerateConfigAndWriteHCL asks questions of the user to work out what the config should be and then writes it
	// using the hclwrite package
	GenerateConfigAndWriteHCL(questioner questions.Questioner, hclBody *hclwrite.Body) error
}

// GetCatalogName returns the name of the plugin as it appears in the plugin catalog. This does not include the plugin
// version, to allow the plugin to be updated without needed to remount the associated instances. However it does
// include the mount path to allow the version of the plugin to vary independently between different mounted instances
// of it.
func (p *PluginConfig) GetCatalogName() string {
	return fmt.Sprintf("%s-%s", p.Type, p.MountPath)
}

// GetFileName returns the filename of the plugin as it will be found in the plugin directory. This includes the plugin
// version to allow different versions of the plugin to be present on the Vault server, and for the catalog entries to
// reference different ones depending on their mounts' use case. Alternatively, if PluginConfig.Filename is specified, then
// the default behaviour will be overridden and it will be returned instead.
func (p *PluginConfig) GetFileName() string {
	if p.Filename != "" {
		return p.Filename
	}

	return fmt.Sprintf("%s_%s", p.Type, p.Version)
}

// WriteHCL uses the hclwrite package to encode itself into HCL
func (p *PluginConfig) WriteHCL(hclBody *hclwrite.Body) {
	hclBody.AppendNewline()
	pluginConfigBlock := hclBody.AppendNewBlock("plugin", []string{p.Type, p.MountPath})
	pluginConfigBody := pluginConfigBlock.Body()

	pluginConfigBody.SetAttributeValue("version", cty.StringVal(p.Version))
}
