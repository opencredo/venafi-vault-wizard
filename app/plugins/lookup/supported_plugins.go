package lookup

import (
	"errors"

	"github.com/opencredo/venafi-vault-wizard/app/plugins"
	pki_backend "github.com/opencredo/venafi-vault-wizard/app/plugins/venafi/pki-backend"
	pki_monitor "github.com/opencredo/venafi-vault-wizard/app/plugins/venafi/pki-monitor"
)

var ErrPluginNotFound = errors.New("plugin not found")

type PluginImplConstructor func() plugins.PluginImpl

// Map of "constructors" which maps all the supported plugin types to their associated PluginImpl implementations.
// Uses a function to force each instance to be a copy, and to allow injection of fields from the wrapper Plugin struct
var supportedPlugins = map[string]PluginImplConstructor{
	"venafi-pki-backend": func() plugins.PluginImpl {
		return &pki_backend.VenafiPKIBackendConfig{}
	},
	"venafi-pki-monitor": func() plugins.PluginImpl {
		return &pki_monitor.VenafiPKIMonitorConfig{}
	},
}

// GetPlugin goes from a generic config.Plugin and looks up its specific PluginImpl based on the Type field.
func GetPlugin(pluginType string) (plugins.PluginImpl, error) {
	constructor, ok := supportedPlugins[pluginType]
	if !ok {
		return nil, ErrPluginNotFound
	}

	return constructor(), nil
}

// SupportedPluginNames returns a list of the names of all of the supported plugins
func SupportedPluginNames() (names []string) {
	for name := range supportedPlugins {
		names = append(names, name)
	}
	return
}
