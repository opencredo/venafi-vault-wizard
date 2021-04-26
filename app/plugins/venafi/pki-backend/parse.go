package pki_backend

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"

	"github.com/opencredo/venafi-vault-wizard/app/plugins"
)

func ParseVenafiPKIBackendConfig(config *plugins.Plugin, evalContext *hcl.EvalContext) (plugins.PluginImpl, error) {
	pluginConfig := &VenafiPKIBackendConfig{
		MountPath: config.MountPath,
		Version:   config.Version,
	}

	diagnostics := gohcl.DecodeBody(config.Config, evalContext, pluginConfig)
	if diagnostics.HasErrors() {
		return nil, diagnostics
	}

	return pluginConfig, nil
}
