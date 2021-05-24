package commands

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/manifoldco/promptui"
	"github.com/opencredo/venafi-vault-wizard/app/config"
)

func GenerateConfig(configFilePath string) {
	vaultConfig, err := generateVaultConfig()
	if err != nil {
		fmt.Printf("Error while generating Vault config: %v\n", err)
		return
	}

	configuration := hclwrite.NewEmptyFile()
	rootBody := configuration.Body()

	vaultConfig.WriteHCL(rootBody)

	fmt.Println("Config file would be saved to", configFilePath, "with contents\n", string(configuration.Bytes()))
}

func generateVaultConfig() (*config.VaultConfig, error) {
	apiAddressPrompt := promptui.Prompt{
		Label:       "What is Vault's API address?",
		Default:     "http://localhost:8200",
		HideEntered: true,
		AllowEdit:   true,
		Validate: func(input string) error {
			if strings.HasPrefix(input, "$") {
				return nil
			}

			_, err := url.ParseRequestURI(input)
			return err
		},
	}
	apiAddress, err := apiAddressPrompt.Run()
	if err != nil {
		return nil, err
	}

	tokenPrompt := promptui.Prompt{
		Label:       "What token should be used to authenticate with Vault?",
		HideEntered: true,
		Default:     "$VAULT_TOKEN",
	}
	token, err := tokenPrompt.Run()
	if err != nil {
		return nil, err
	}

	containerOrVMPrompt := promptui.Select{
		Label:        "Is Vault running in a VM or a container",
		HideSelected: true,
		Items:        []string{"VM", "Container"},
	}

	_, containerOrVM, err := containerOrVMPrompt.Run()
	if err != nil {
		return nil, err
	}

	var vaultConfig = &config.VaultConfig{
		VaultAddress: apiAddress,
		VaultToken:   token,
	}

	if containerOrVM == "VM" {
		sshPrompt := promptui.Select{
			Label:        "Do you have SSH access to the Vault server(s)",
			HideSelected: true,
			Items:        []string{"Yes", "No"},
		}
		_, useSSH, err := sshPrompt.Run()
		if err != nil {
			return nil, err
		}

		if useSSH == "Yes" {
			// TODO: ask for ssh details
			vaultConfig.SSHConfig = append(vaultConfig.SSHConfig, config.SSH{
				Hostname: "ssh.hostname",
				Username: "vagrant",
				Password: "vagrant",
				Port:     22,
			})

			return vaultConfig, nil
		}
	}

	pluginsIncludedInImagePrompt := promptui.Select{
		Label:        "Are the plugin binaries already included in the server's image",
		HideSelected: true,
		Items:        []string{"Yes", "No"},
	}
	_, pluginsIncludedInImage, err := pluginsIncludedInImagePrompt.Run()
	if err != nil {
		return nil, err
	}

	if pluginsIncludedInImage == "No" {
		return nil, fmt.Errorf("you must either have SSH access, or include the plugin binaries externally")
	}

	return vaultConfig, nil
}
