package lib

import (
	"fmt"

	vaultAPI "github.com/hashicorp/vault/api"
)

// VaultAPIWrapper encapsulates the dependency on the HashiCorp Go Vault package, both to allow it to be injected, but
// also to make its interface slightly simpler to its clients
type VaultAPIWrapper interface {
	SetAddress(address string) error
	SetToken(token string)
	Read(path string) (map[string]interface{}, error)
	Write(path string, data map[string]interface{}) (map[string]interface{}, error)
	RegisterPlugin(input *vaultAPI.RegisterPluginInput) error
	GetPlugin(input *vaultAPI.GetPluginInput) (*vaultAPI.GetPluginResponse, error)
	ReloadPlugin(input *vaultAPI.ReloadPluginInput) (string, error)
	Mount(path string, input *vaultAPI.MountInput) error
	ListMounts() (map[string]*vaultAPI.MountOutput, error)
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

	if secret != nil {
		return secret.Data, nil
	} else {
		return nil, fmt.Errorf("no data found at path %s", path)
	}
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

func (v *vaultAPIClient) GetPlugin(input *vaultAPI.GetPluginInput) (*vaultAPI.GetPluginResponse, error) {
	plugin, err := v.Sys().GetPlugin(input)
	return plugin, normaliseError(err)
}

func (v *vaultAPIClient) ReloadPlugin(input *vaultAPI.ReloadPluginInput) (string, error) {
	reloadID, err := v.Sys().ReloadPlugin(input)
	return reloadID, normaliseError(err)
}

func (v *vaultAPIClient) Mount(path string, input *vaultAPI.MountInput) error {
	err := v.Sys().Mount(path, input)
	return normaliseError(err)
}

func (v *vaultAPIClient) ListMounts() (map[string]*vaultAPI.MountOutput, error) {
	mounts, err := v.Sys().ListMounts()
	return mounts, normaliseError(err)
}
