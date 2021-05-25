package commands

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/opencredo/venafi-vault-wizard/app/config"
	"github.com/opencredo/venafi-vault-wizard/app/plugins/lookup"
	"github.com/opencredo/venafi-vault-wizard/app/questions"
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
		return
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
	answers := questions.NewAnswerQueue()
	err := questions.AskQuestions([]questions.Question{
		&questions.OpenEndedQuestion{
			Question:  "What is Vault's API address?",
			Default:   "http://localhost:8200",
			AllowEdit: true,
			Validate: func(input string) error {
				if strings.HasPrefix(input, "$") {
					return nil
				}

				_, err := url.ParseRequestURI(input)
				return err
			},
		},
		&questions.OpenEndedQuestion{
			Question: "What token should be used to authenticate with Vault?",
			Default:  "$VAULT_TOKEN",
		},
		&questions.QuestionBranch{
			ConditionQuestion: &questions.ClosedQuestion{
				Question: "Is Vault running in a VM or a container",
				Items:    []string{"VM", "Container"},
			},
			ConditionAnswer: "VM",
			BranchA: []questions.Question{
				&questions.ClosedQuestion{
					Question: "Do you have SSH access to the Vault server(s)",
					Items:    []string{"Yes", "No"},
				},
			},
			BranchB: []questions.Question{
				&questions.ClosedQuestion{
					Question: "Are the plugin binaries already included in the server's image",
					Items:    []string{"Yes", "No"},
				},
			},
		},
	}, answers)
	if err != nil {
		return nil, err
	}

	var vaultConfig = &config.VaultConfig{
		VaultAddress: string(*answers.Pop()),
		VaultToken:   string(*answers.Pop()),
	}

	var containerOrVM = string(*answers.Pop())
	var pluginIncludedInImage questions.Answer
	if containerOrVM == "VM" {
		var useSSH = string(*answers.Pop())
		if useSSH == "Yes" {
			sshConfigs, err := generateSSHConfigs()
			if err != nil {
				return nil, err
			}
			vaultConfig.SSHConfig = sshConfigs

			return vaultConfig, nil
		}

		pluginIncludedInImage, err = questions.AskSingleQuestion(&questions.ClosedQuestion{
			Question: "Are the plugin binaries already included in the server's image",
			Items:    []string{"Yes", "No"},
		})
		if err != nil {
			return nil, err
		}
	} else {
		pluginIncludedInImage = *answers.Pop()
	}

	if pluginIncludedInImage == "No" {
		return nil, fmt.Errorf("you must either have SSH access, or include the plugin binaries externally")
	}

	return vaultConfig, nil
}

func generatePluginsConfig() ([]*hclwrite.Block, error) {
	var pluginBlocks []*hclwrite.Block
	for i := 1; true; i++ {
		answers := questions.NewAnswerQueue()
		err := questions.AskQuestions([]questions.Question{
			&questions.ClosedQuestion{
				Question: "Which plugin would you like to configure",
				Items:    lookup.SupportedPluginNames(),
			},
			&questions.OpenEndedQuestion{
				Question: "Which version of the plugin would you like to use?",
			},
			&questions.OpenEndedQuestion{
				Question: "Which Vault path should the plugin be mounted at?",
			},
		}, answers)
		if err != nil {
			return nil, err
		}
		pluginType, version, mountPath := string(*answers.Pop()), string(*answers.Pop()), string(*answers.Pop())

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

		morePlugins, err := questions.AskSingleQuestion(&questions.ClosedQuestion{
			Question: fmt.Sprintf("You have configured %d plugins, are there more", i),
			Items:    []string{"Yes", "No that's it"},
		})
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
	ha, err := questions.AskSingleQuestion(&questions.ClosedQuestion{
		Question: "Is Vault running in High-Availability (HA) mode",
		Items:    []string{"Yes", "No, just one node"},
	})
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

			moreSSHs, err := questions.AskSingleQuestion(&questions.ClosedQuestion{
				Question: fmt.Sprintf("You have configured %d Vault replicas, are there more", i),
				Items:    []string{"Yes", "No, that's it"},
			})
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
	answers := questions.NewAnswerQueue()
	err := questions.AskQuestions([]questions.Question{
		&questions.OpenEndedQuestion{
			Question: "What is the hostname of the Vault server?",
		},
		&questions.OpenEndedQuestion{
			Question: "What is the SSH username to log into the Vault server?",
		},
		&questions.OpenEndedQuestion{
			Question: "What is the SSH password to log into the Vault server?",
		},
		&questions.OpenEndedQuestion{
			Question: "What is the SSH port for logging into the Vault server?",
			Default:  "22",
			Validate: func(input string) error {
				_, err := strconv.ParseUint(input, 10, 16)
				if err != nil {
					return fmt.Errorf("SSH port must be an integer")
				}
				return nil
			},
		},
	}, answers)
	if err != nil {
		return nil, err
	}

	sshConfig := &config.SSH{
		Hostname: string(*answers.Pop()),
		Username: string(*answers.Pop()),
		Password: string(*answers.Pop()),
	}
	sshPort, _ := strconv.ParseUint(string(*answers.Pop()), 10, 16)
	sshConfig.Port = uint(sshPort)

	return sshConfig, nil
}
