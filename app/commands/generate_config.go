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

func GenerateConfig(configFilePath string, questioner questions.Questioner) {
	file, err := os.OpenFile(configFilePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("Error opening file %s for writing: %s\n", configFilePath, err)
		return
	}
	defer file.Close()

	vaultConfig, err := generateVaultConfig(questioner)
	if err != nil {
		fmt.Printf("Error while generating Vault config: %v\n", err)
		return
	}

	pluginBlocks, err := generatePluginsConfig(questioner)
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

func generateVaultConfig(questioner questions.Questioner) (*config.VaultConfig, error) {
	q := map[string]questions.Question{
		"api_addr": questioner.NewOpenEndedQuestion(&questions.OpenEndedQuestion{
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
		}),
		"token": questioner.NewOpenEndedQuestion(&questions.OpenEndedQuestion{
			Question: "What token should be used to authenticate with Vault?",
			Default:  "$VAULT_TOKEN",
		}),
		"vm/container": questioner.NewClosedQuestion(&questions.ClosedQuestion{
			Question: "Is Vault running in a VM or a container",
			Items:    []string{"VM", "Container"},
		}),
		"ssh": questioner.NewClosedQuestion(&questions.ClosedQuestion{
			Question: "Do you have SSH access to the Vault server(s)",
			Items:    []string{"Yes", "No"},
		}),
		"binaries_incl": questioner.NewClosedQuestion(&questions.ClosedQuestion{
			Question: "Are the plugin binaries already included in the server's image",
			Items:    []string{"Yes", "No"},
		}),
	}
	err := questions.AskQuestions([]questions.Question{
		q["api_addr"],
		q["token"],
		&questions.QuestionBranch{
			ConditionQuestion: q["vm/container"],
			ConditionAnswer:   "VM",
			BranchA: []questions.Question{
				&questions.QuestionBranch{
					ConditionQuestion: q["ssh"],
					ConditionAnswer:   "No",
					BranchA: []questions.Question{
						q["binaries_incl"],
					},
				},
			},
			BranchB: []questions.Question{
				q["binaries_incl"],
			},
		},
	})
	if err != nil {
		return nil, err
	}

	var vaultConfig = &config.VaultConfig{
		VaultAddress: string(q["api_addr"].Answer()),
		VaultToken:   string(q["token"].Answer()),
	}

	if q["vm/container"].Answer() == "VM" && q["ssh"].Answer() == "Yes" {
		sshConfigs, err := generateSSHConfigs(questioner)
		if err != nil {
			return nil, err
		}
		vaultConfig.SSHConfig = sshConfigs
	} else if q["binaries_incl"].Answer() == "No" {
		return nil, fmt.Errorf("you must either have SSH access, or include the plugin binaries externally")
	}

	return vaultConfig, nil
}

func generatePluginsConfig(questioner questions.Questioner) ([]*hclwrite.Block, error) {
	var pluginBlocks []*hclwrite.Block
	var buildArch string

	for i := 1; true; i++ {
		q := map[string]questions.Question{
			"type": questioner.NewClosedQuestion(&questions.ClosedQuestion{
				Question: "Which plugin would you like to configure",
				Items:    lookup.SupportedPluginNames(),
			}),
			"version": questioner.NewOpenEndedQuestion(&questions.OpenEndedQuestion{
				Question: "Which version of the plugin would you like to use?",
			}),
			"mount_path": questioner.NewOpenEndedQuestion(&questions.OpenEndedQuestion{
				Question: "Which Vault path should the plugin be mounted at?",
			}),
			"os_type": questioner.NewClosedQuestion(&questions.ClosedQuestion{
				Question: "Which OS type is Vault running on?",
				Items: []string{"Linux", "Mac OS"},
			}),
			"os_arch": questioner.NewClosedQuestion(&questions.ClosedQuestion{
				Question: "Which OS architecture is used?",
				Items: []string{"64bit", "32bit"},
			}),
		}
		err := questions.AskQuestions([]questions.Question{
			q["type"],
			q["version"],
			q["mount_path"],
			&questions.QuestionBranch{
				ConditionQuestion: q["os_type"],
				ConditionAnswer: "Linux",
				BranchA: []questions.Question{
					q["os_arch"],
				},
			},
		})
		if err != nil {
			return nil, err
		}
		pluginType, version, mountPath := string(q["type"].Answer()), string(q["version"].Answer()), string(q["mount_path"].Answer())

		switch osType := string(q["os_type"].Answer()); osType {
		case "Mac OS":
			buildArch = "darwin"
		case "Linux":
			buildArch = "linux"
			if string(q["os_arch"].Answer()) == "32bit" {
				buildArch = buildArch + "86"
			}
		}
		pluginImpl, err := lookup.GetPlugin(pluginType)
		if err != nil {
			return nil, err
		}

		pluginBlock := hclwrite.NewBlock("plugin", []string{pluginType, mountPath})
		pluginBody := pluginBlock.Body()
		pluginBody.SetAttributeValue("version", cty.StringVal(version))
		pluginBody.SetAttributeValue("build_arch", cty.StringVal(buildArch))
		err = pluginImpl.GenerateConfigAndWriteHCL(questioner, pluginBody)
		if err != nil {
			return nil, err
		}

		pluginBlocks = append(pluginBlocks, pluginBlock)

		morePluginsQuestion := questioner.NewClosedQuestion(&questions.ClosedQuestion{
			Question: fmt.Sprintf("You have configured %d plugins, are there more", i),
			Items:    []string{"Yes", "No that's it"},
		})
		err = morePluginsQuestion.Ask()
		if err != nil {
			return nil, err
		}

		if morePluginsQuestion.Answer() != "Yes" {
			break
		}
	}
	return pluginBlocks, nil
}

func generateSSHConfigs(questioner questions.Questioner) ([]config.SSH, error) {
	haQuestion := questioner.NewClosedQuestion(&questions.ClosedQuestion{
		Question: "Is Vault running in High-Availability (HA) mode",
		Items:    []string{"Yes", "No, just one node"},
	})
	err := haQuestion.Ask()
	if err != nil {
		return nil, err
	}

	var sshConfigs []config.SSH

	if haQuestion.Answer() == "Yes" {
		for i := 1; true; i++ {
			sshConfig, err := generateSSHConfig(questioner)
			if err != nil {
				return nil, err
			}
			sshConfigs = append(sshConfigs, *sshConfig)

			moreSSHsQuestion := questioner.NewClosedQuestion(&questions.ClosedQuestion{
				Question: fmt.Sprintf("You have configured %d Vault replicas, are there more", i),
				Items:    []string{"Yes", "No, that's it"},
			})
			err = moreSSHsQuestion.Ask()
			if err != nil {
				return nil, err
			}

			if moreSSHsQuestion.Answer() != "Yes" {
				break
			}
		}
	} else {
		sshConfig, err := generateSSHConfig(questioner)
		if err != nil {
			return nil, err
		}
		sshConfigs = append(sshConfigs, *sshConfig)
	}

	return sshConfigs, nil
}

func generateSSHConfig(questioner questions.Questioner) (*config.SSH, error) {
	q := map[string]questions.Question{
		"hostname": questioner.NewOpenEndedQuestion(&questions.OpenEndedQuestion{
			Question: "What is the hostname of the Vault server?",
		}),
		"username": questioner.NewOpenEndedQuestion(&questions.OpenEndedQuestion{
			Question: "What is the SSH username to log into the Vault server?",
		}),
		"password": questioner.NewOpenEndedQuestion(&questions.OpenEndedQuestion{
			Question: "What is the SSH password to log into the Vault server?",
		}),
		"port": questioner.NewOpenEndedQuestion(&questions.OpenEndedQuestion{
			Question: "What is the SSH port for logging into the Vault server?",
			Default:  "22",
			Validate: func(input string) error {
				_, err := strconv.ParseUint(input, 10, 16)
				if err != nil {
					return fmt.Errorf("SSH port must be an integer")
				}
				return nil
			},
		}),
	}
	err := questions.AskQuestions([]questions.Question{
		q["hostname"],
		q["username"],
		q["password"],
		q["port"],
	})
	if err != nil {
		return nil, err
	}

	sshConfig := &config.SSH{
		Hostname: string(q["hostname"].Answer()),
		Username: string(q["username"].Answer()),
		Password: string(q["password"].Answer()),
	}
	sshPort, _ := strconv.ParseUint(string(q["port"].Answer()), 10, 16)
	sshConfig.Port = uint(sshPort)

	return sshConfig, nil
}
