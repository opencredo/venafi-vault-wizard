package pki_monitor

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"

	"github.com/opencredo/venafi-vault-wizard/app/plugins"
)

func (c *VenafiPKIMonitorConfig) ParseConfig(config *plugins.PluginConfig, evalContext *hcl.EvalContext) error {
	c.MountPath = config.MountPath
	c.Version = config.Version
	c.BuildArch = config.BuildArch

	diagnostics := gohcl.DecodeBody(config.Config, evalContext, c)
	if diagnostics.HasErrors() {
		return diagnostics
	}

	return nil
}
