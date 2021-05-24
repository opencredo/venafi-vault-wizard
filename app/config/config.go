package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"

	"github.com/opencredo/venafi-vault-wizard/app/plugins"
	"github.com/opencredo/venafi-vault-wizard/app/plugins/lookup"
)

type Config struct {
	Vault   VaultConfig      `hcl:"vault,block"`
	Plugins []plugins.Plugin `hcl:"plugin,block"`
}

// NewConfig decodes an HCL configuration file into a Config struct, returning an error upon failure. It takes filename
// as a parameter to use in error messages while parsing, and a byte slice containing the actual configuration itself.
// It then uses hclsimple to parse the configuration and validates it using Config.Validate() before returning it.
// If this function returns without an error then the config should be valid to use.
func NewConfig(filename string, src []byte) (*Config, error) {
	config := new(Config)

	// Add a custom env function to the context so that config files can use it
	configEvaluationContext := &hcl.EvalContext{
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
						return cty.NullVal(cty.String), fmt.Errorf("environment variable not set")
					}

					return cty.StringVal(env), nil
				},
			}),
		},
	}

	err := hclsimple.Decode(filename, src, configEvaluationContext, config)
	if err != nil {
		return nil, err
	}

	err = config.Validate()
	if err != nil {
		return nil, err
	}

	for i, plugin := range config.Plugins {
		pluginImpl, err := lookup.LookupPlugin(&plugin, configEvaluationContext)
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

// NewConfigFromFile reads the file from filename and calls NewConfig with its contents
func NewConfigFromFile(filename string) (*Config, error) {
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("error config file %s not found: %w", filename, err)
		}

		return nil, fmt.Errorf("can't read %s: %w", filename, err)
	}

	return NewConfig(filename, src)
}

func (c *Config) Validate() error {
	err := c.Vault.Validate()
	if err != nil {
		return err
	}

	return nil
}

func WriteStringAttributeToHCL(attributeName, input string, body *hclwrite.Body) {
	if strings.HasPrefix(input, "$") {
		body.SetAttributeRaw(attributeName, envFunctionTokens(input[1:]))
		return
	}

	body.SetAttributeValue(attributeName, cty.StringVal(input))
}

func envFunctionTokens(environmentVariableName string) hclwrite.Tokens {
	return hclwrite.Tokens{
		&hclwrite.Token{
			Type:  hclsyntax.TokenIdent,
			Bytes: []byte("env"),
		},
		&hclwrite.Token{
			Type:  hclsyntax.TokenOParen,
			Bytes: []byte("("),
		},
		&hclwrite.Token{
			Type:  hclsyntax.TokenQuotedLit,
			Bytes: []byte("\"" + environmentVariableName + "\""),
		},
		&hclwrite.Token{
			Type:  hclsyntax.TokenCParen,
			Bytes: []byte(")"),
		},
	}
}
