package lib

import vaultAPI "github.com/hashicorp/vault/api"

// VaultAPIWrapper encapsulates the dependency on the HashiCorp Go Vault package, both to allow it to be injected, but
// also to make its interface slightly simpler to its clients
type VaultAPIWrapper interface {
	SetAddress(address string) error
	SetToken(token string)
	Read(path string) (map[string]interface{}, error)
	Write(path string, data map[string]interface{}) (map[string]interface{}, error)
	RegisterPlugin(input *vaultAPI.RegisterPluginInput) error
	Mount(path string, input *vaultAPI.MountInput) error
}

type vaultAPIClient struct {
	*vaultAPI.Client
}

func NewVaultAPI() VaultAPIWrapper {
	client, _ := vaultAPI.NewClient(vaultAPI.DefaultConfig())
	return &vaultAPIClient{client}
}

func (v *vaultAPIClient) Read(path string) (map[string]interface{}, error) {
	secret, err := v.Logical().Read(path)
	if err != nil {
		return nil, normaliseError(err)
	}

	return secret.Data, nil
}

func (v *vaultAPIClient) Write(path string, data map[string]interface{}) (map[string]interface{}, error) {
	secret, err := v.Logical().Write(path, data)
	if err != nil {
		return nil, normaliseError(err)
	}

	if secret == nil {
		return nil, nil
	} else {
		return secret.Data, nil
	}
}

func (v *vaultAPIClient) RegisterPlugin(input *vaultAPI.RegisterPluginInput) error {
	err := v.Sys().RegisterPlugin(input)
	return normaliseError(err)
}

func (v *vaultAPIClient) Mount(path string, input *vaultAPI.MountInput) error {
	err := v.Sys().Mount(path, input)
	return normaliseError(err)
}
