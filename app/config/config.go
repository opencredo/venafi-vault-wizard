package config

import (
	"fmt"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/opencredo/venafi-vault-wizard/app/vault/api"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"

	"github.com/opencredo/venafi-vault-wizard/app/plugins"
	"github.com/opencredo/venafi-vault-wizard/app/plugins/lookup"
)

// configStage1 is a private struct used to decode only the Vault API parameters. This allows the next stage of decoding
// to retrieve values from the Vault API
type configStage1 struct {
	Vault struct {
		Address string   `hcl:"api_address"`
		Token   string   `hcl:"token"`
		Rest    hcl.Body `hcl:",remain"`
	} `hcl:"vault,block"`
	Rest hcl.Body `hcl:",remain"`
}

type Config struct {
	Vault   VaultConfig            `hcl:"vault,block"`
	Plugins []plugins.PluginConfig `hcl:"plugin,block"`
}

type ConfigParser struct {
	fileName     string
	fileContents []byte
	evalContext  *hcl.EvalContext
}

// NewConfigParser returns a ConfigParser which, initially, only supports an `env()` function in the config. Calling
// ConfigParser.SetVaultClient afterwards adds support for the `secret()` function too, to allow parameters to be
// retrieved from Vault. The ConfigParser.GetVaultAPIDetails method can be used to partially decode the config to get
// just the API address and token, which is sufficient for creating an api.VaultAPIClient to give to
// ConfigParser.SetVaultClient. This means that the rest of the config can be fully decoded afterwards using
// ConfigParser.GetConfig, and the rest of the config can take advantage of the `secret()` function. This function
// should be provided with the config file name (for constructing the diagnostic messages), and a byte slice containing
// the contents of the config file.
func NewConfigParser(filename string, src []byte) *ConfigParser {
	return &ConfigParser{
		fileName:     filename,
		fileContents: src,
		evalContext: &hcl.EvalContext{
			Functions: map[string]function.Function{
				"env": function.New(&function.Spec{
					Params: []function.Parameter{
						{
							Name:             "name",
							Type:             cty.String,
							AllowDynamicType: true,
						},
					},
					Type: function.StaticReturnType(cty.String),
					Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
						env := os.Getenv(args[0].AsString())
						if env == "" {
							return cty.NilVal, fmt.Errorf("environment variable not set")
						}

						return cty.StringVal(env), nil
					},
				}),
			},
		},
	}
}

// GetVaultAPIDetails partially decodes the config in order to retrieve an API address and Token for Vault.
func (p *ConfigParser) GetVaultAPIDetails() (apiAddr, token string, err error) {
	config := new(configStage1)

	err = hclsimple.Decode(p.fileName, p.fileContents, p.evalContext, config)
	if err != nil {
		return "", "", err
	}

	return config.Vault.Address, config.Vault.Token, nil
}

// SetVaultClient provides the config parser with access to the Vault API which allows the config to use the `secret()`
// function to retrieve Vault secrets automatically. Unless this method is called, ConfigParser.GetConfig will fail if
// the config tries to use the `secret()` function.
func (p *ConfigParser) SetVaultClient(vaultClient api.VaultAPIClient) {
	p.evalContext.Functions["secret"] = function.New(&function.Spec{
		Params: []function.Parameter{
			{
				Name:             "vaultPath",
				Type:             cty.String,
				AllowDynamicType: true,
			},
			{
				Name:             "fieldKey",
				Type:             cty.String,
				AllowDynamicType: true,
			},
		},
		Type: function.StaticReturnType(cty.String),
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			path := args[0].AsString()
			field := args[1].AsString()

			value, err := vaultClient.ReadValue(path)
			if err != nil {
				return cty.NilVal, fmt.Errorf("error reading Vault secret: %w", err)
			}

			return cty.StringVal(fmt.Sprintf("%v", value[field])), nil
		},
	})
}

// GetConfig decodes an HCL configuration file into a Config struct, returning an error upon failure.
// It uses hclsimple to parse the configuration and validates it using Config.Validate() before returning it.
// If this function returns without an error then the config should be valid to use.
func (p *ConfigParser) GetConfig() (*Config, error) {
	config := new(Config)

	err := hclsimple.Decode(p.fileName, p.fileContents, p.evalContext, config)
	if err != nil {
		return nil, err
	}

	err = config.Validate()
	if err != nil {
		return nil, err
	}

	for i, plugin := range config.Plugins {
		pluginImpl, err := lookup.GetPlugin(plugin.Type)
		if err != nil {
			return nil, err
		}

		err = pluginImpl.ParseConfig(&plugin, p.evalContext)
		if err != nil {
			return nil, err
		}

		err = pluginImpl.ValidateConfig()
		if err != nil {
			return nil, err
		}

		config.Plugins[i].Impl = pluginImpl
	}

	return config, nil
}

func (c *Config) Validate() error {
	err := c.Vault.Validate()
	if err != nil {
		return err
	}

	return nil
}
