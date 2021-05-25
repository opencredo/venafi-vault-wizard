package commands

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/manifoldco/promptui"
	"github.com/opencredo/venafi-vault-wizard/app/config"
	"github.com/opencredo/venafi-vault-wizard/app/plugins/lookup"
	"github.com/zclconf/go-cty/cty"
)

func GenerateConfig(configFilePath string) {
	file, err := os.OpenFile(configFilePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("Error opening file %s for writing: %s\n", configFilePath, err)
		return
	}
	defer file.Close()

	vaultConfig, err := generateVaultConfig()
	if err != nil {
		fmt.Printf("Error while generating Vault config: %v\n", err)
		return
	}

	pluginBlocks, err := generatePluginsConfig()
	if err != nil {
		fmt.Printf("Error while generating plugins config: %v\n", err)
	}

	configuration := hclwrite.NewEmptyFile()
	rootBody := configuration.Body()

	vaultConfig.WriteHCL(rootBody)
	for _, pluginBlock := range pluginBlocks {
		rootBody.AppendNewline()
		rootBody.AppendBlock(pluginBlock)
	}

	_, err = configuration.WriteTo(file)
	if err != nil {
		fmt.Printf("Error writing config to %s: %s\n", configFilePath, err)
		return
	}

	fmt.Printf("Config successfully written to %s\n", configFilePath)
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
			sshConfigs, err := generateSSHConfigs()
			if err != nil {
				return nil, err
			}
			vaultConfig.SSHConfig = sshConfigs

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

func generatePluginsConfig() ([]*hclwrite.Block, error) {
	var pluginBlocks []*hclwrite.Block
	for i := 1; true; i++ {
		pluginTypePrompt := promptui.Select{
			Label:        "Which plugin would you like to configure",
			Items:        lookup.SupportedPluginNames(),
			HideSelected: true,
		}
		_, pluginType, err := pluginTypePrompt.Run()
		if err != nil {
			return nil, err
		}

		versionPrompt := promptui.Prompt{
			Label:       "Which version of the plugin would you like to use?",
			HideEntered: true,
		}
		version, err := versionPrompt.Run()
		if err != nil {
			return nil, err
		}

		mountPathPrompt := promptui.Prompt{
			Label:       "Which Vault path should the plugin be mounted at?",
			HideEntered: true,
		}
		mountPath, err := mountPathPrompt.Run()
		if err != nil {
			return nil, err
		}

		pluginImpl, err := lookup.GetPlugin(pluginType)
		if err != nil {
			return nil, err
		}

		pluginBlock := hclwrite.NewBlock("plugin", []string{pluginType, mountPath})
		pluginBody := pluginBlock.Body()
		pluginBody.SetAttributeValue("version", cty.StringVal(version))
		err = pluginImpl.GenerateConfigAndWriteHCL(pluginBody)
		if err != nil {
			return nil, err
		}

		pluginBlocks = append(pluginBlocks, pluginBlock)

		morePluginsPrompt := promptui.Select{
			Label:        fmt.Sprintf("You have configured %d plugins, are there more", i),
			HideSelected: true,
			Items:        []string{"Yes", "No that's it"},
		}
		_, morePlugins, err := morePluginsPrompt.Run()
		if err != nil {
			return nil, err
		}

		if morePlugins != "Yes" {
			break
		}
	}
	return pluginBlocks, nil
}

func generateSSHConfigs() ([]config.SSH, error) {
	haPrompt := promptui.Select{
		Label:        "Is Vault running in High-Availability (HA) mode",
		HideSelected: true,
		Items:        []string{"Yes", "No, just one node"},
	}
	_, ha, err := haPrompt.Run()
	if err != nil {
		return nil, err
	}

	var sshConfigs []config.SSH

	if ha == "Yes" {
		for i := 1; true; i++ {
			sshConfig, err := generateSSHConfig()
			if err != nil {
				return nil, err
			}
			sshConfigs = append(sshConfigs, *sshConfig)

			moreSSHsPrompt := promptui.Select{
				Label:        fmt.Sprintf("You have configured %d Vault replicas, are there more", i),
				HideSelected: true,
				Items:        []string{"Yes", "No, that's it"},
			}
			_, moreSSHs, err := moreSSHsPrompt.Run()
			if err != nil {
				return nil, err
			}

			if moreSSHs != "Yes" {
				break
			}
		}
	} else {
		sshConfig, err := generateSSHConfig()
		if err != nil {
			return nil, err
		}
		sshConfigs = append(sshConfigs, *sshConfig)
	}

	return sshConfigs, nil
}

func generateSSHConfig() (*config.SSH, error) {
	hostnamePrompt := promptui.Prompt{
		Label:       "What is the hostname of the Vault server?",
		HideEntered: true,
	}
	hostname, err := hostnamePrompt.Run()
	if err != nil {
		return nil, err
	}

	usernamePrompt := promptui.Prompt{
		Label:       "What is the SSH username to log into the Vault server?",
		HideEntered: true,
	}
	username, err := usernamePrompt.Run()
	if err != nil {
		return nil, err
	}

	passwordPrompt := promptui.Prompt{
		Label:       "What is the SSH password to log into the Vault server?",
		HideEntered: true,
	}
	password, err := passwordPrompt.Run()
	if err != nil {
		return nil, err
	}

	portPrompt := promptui.Prompt{
		Label:       "What is the SSH port for logging into the Vault server?",
		Default:     "22",
		HideEntered: true,
		Validate: func(input string) error {
			_, err := strconv.Atoi(input)
			if err != nil {
				return fmt.Errorf("SSH port must be an integer")
			}
			return nil
		},
	}
	portString, err := portPrompt.Run()
	if err != nil {
		return nil, err
	}
	port, _ := strconv.Atoi(portString)

	return &config.SSH{
		Hostname: hostname,
		Username: username,
		Password: password,
		Port:     uint(port),
	}, nil
}
