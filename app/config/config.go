package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

var ErrBlankParam = errors.New("cannot be blank")
var ErrConflictingBlocks = errors.New("one of the blocks must be defined")

type Config struct {
	Vault      VaultConfig             `hcl:"vault,block"`
	PKIBackend *VenafiPKIBackendConfig `hcl:"venafi_pki_backend,block"`
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

	err = c.PKIBackend.Validate()
	if err != nil {
		return err
	}

	return nil
}
