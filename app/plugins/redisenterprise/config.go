package redisenterprise

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/opencredo/venafi-vault-wizard/app/config/errors"
	"github.com/opencredo/venafi-vault-wizard/app/plugins"
	"github.com/opencredo/venafi-vault-wizard/app/questions"
)

type RedisEnterpriseConfig struct {
	// MountPath is not decoded directly by using the struct tags, and is instead populated by
	// ParseConfig when it is initialised
	MountPath string
	// Version is not decoded directly by using the struct tags, and is instead populated by ParseConfig
	// when it is initialised
	Version string
	// PluginName is not decoded directly by using the struct tags, and is instead populated by ParseConfig
	// when it is initialised
	PluginName string

	Roles []Role `hcl:"role,block"`
}

type Role struct {
	Name   string `hcl:"role,label"`
	DBRole string `hcl:"db_role"`
	ACL    string `hcl:"acl,optional"`

	DBCluster Cluster `hcl:"db_cluster,block"`
}

type Cluster struct {
	Name         string `hcl:"cluster,label"`
	Username     string `hcl:"username"`
	Password     string `hcl:"password"`
	URL          string `hcl:"url"`
	DatabaseName string `hcl:"db_name,optional"`
}

func (c *RedisEnterpriseConfig) ParseConfig(config *plugins.PluginConfig, evalContext *hcl.EvalContext) error {
	c.MountPath = config.MountPath
	c.Version = config.Version
	c.PluginName = config.GetCatalogName()

	diagnostics := gohcl.DecodeBody(config.Config, evalContext, c)
	if diagnostics.HasErrors() {
		return diagnostics
	}

	return nil
}

func (c *RedisEnterpriseConfig) ValidateConfig() error {
	if len(c.Roles) == 0 {
		return fmt.Errorf("error at least one role must be provided: %w", errors.ErrBlankParam)
	}
	for _, role := range c.Roles {
		err := role.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Role) Validate() error {
	// TODO: check valid permutations of options
	return nil
}

func (c *RedisEnterpriseConfig) GenerateConfigAndWriteHCL(questioner questions.Questioner, hclBody *hclwrite.Body) error {
	panic("implement me")
}
