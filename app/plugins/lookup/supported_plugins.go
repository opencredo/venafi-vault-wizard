package lookup

import (
	"errors"

	"github.com/hashicorp/hcl/v2"

	"github.com/opencredo/venafi-vault-wizard/app/plugins"
	pki_backend "github.com/opencredo/venafi-vault-wizard/app/plugins/venafi/pki-backend"
	pki_monitor "github.com/opencredo/venafi-vault-wizard/app/plugins/venafi/pki-monitor"
)

var ErrPluginNotFound = errors.New("plugin not found")

// Map of "constructors" which maps all the supported plugin types to their associated PluginImpl implementations.
// Uses a function to force each instance to be a copy, and to allow injection of fields from the wrapper Plugin struct
type PluginImplConstructor func(*plugins.Plugin, *hcl.EvalContext) (plugins.PluginImpl, error)

var supportedPlugins = map[string]PluginImplConstructor{
	"venafi-pki-backend": pki_backend.ParseVenafiPKIBackendConfig,
	"venafi-pki-monitor": pki_monitor.ParseVenafiPKIMonitorConfig,
}

// LookupPlugin goes from a generic config.Plugin and looks up its specific PluginImpl based on the Type
// field. It takes a pointer to an hcl.EvalContext in order to provide the same config functions available to the rest
// of the config, and in future maybe some global variables, and then decodes the plugin-specific part of the plugin
// block.
func LookupPlugin(config *plugins.Plugin, evalContext *hcl.EvalContext) (plugins.PluginImpl, error) {
	constructor, ok := supportedPlugins[config.Type]
	if !ok {
		return nil, ErrPluginNotFound
	}

	plugin, err := constructor(config, evalContext)
	if err != nil {
		return nil, err
	}

	return plugin, nil
}
