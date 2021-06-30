package lookup

import (
	"errors"

	"github.com/opencredo/venafi-vault-wizard/app/plugins"
	"github.com/opencredo/venafi-vault-wizard/app/plugins/redisenterprise"
	pki_backend "github.com/opencredo/venafi-vault-wizard/app/plugins/venafi/pki-backend"
	pki_monitor "github.com/opencredo/venafi-vault-wizard/app/plugins/venafi/pki-monitor"
	"github.com/opencredo/venafi-vault-wizard/app/vault/api"
)

var ErrPluginNotFound = errors.New("plugin not found")

// PluginConstructor returns the implementation of plugins.Plugin, and the type of Vault plugin that it corresponds to
type PluginConstructor func() (plugins.Plugin, api.PluginType)

// Map of "constructors" which maps all the supported plugin types to their associated PluginImpl implementations.
// Uses a function to force each instance to be a copy, and to allow injection of fields from the wrapper Plugin struct
var supportedPlugins = map[string]PluginConstructor{
	"venafi-pki-backend": func() (plugins.Plugin, api.PluginType) {
		return &pki_backend.VenafiPKIBackendConfig{}, api.PluginTypeSecrets
	},
	"venafi-pki-monitor": func() (plugins.Plugin, api.PluginType) {
		return &pki_monitor.VenafiPKIMonitorConfig{}, api.PluginTypeSecrets
	},
	"redisenterprise": func() (plugins.Plugin, api.PluginType) {
		return &redisenterprise.RedisEnterpriseConfig{}, api.PluginTypeDatabase
	},
}

// GetPlugin goes from a pluginType string and looks up its specific plugins.Plugin implementation and its api.PluginType
func GetPlugin(pluginType string) (plugins.Plugin, api.PluginType, error) {
	constructor, ok := supportedPlugins[pluginType]
	if !ok {
		return nil, 0, ErrPluginNotFound
	}

	i, t := constructor()
	return i, t, nil
}

// SupportedPluginNames returns a list of the names of all of the supported plugins
func SupportedPluginNames() (names []string) {
	for name := range supportedPlugins {
		names = append(names, name)
	}
	return
}
