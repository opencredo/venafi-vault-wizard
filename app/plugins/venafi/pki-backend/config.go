package pki_backend

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/opencredo/venafi-vault-wizard/app/config/errors"
	"github.com/opencredo/venafi-vault-wizard/app/plugins/venafi"
	"github.com/opencredo/venafi-vault-wizard/app/questions"
)

type VenafiPKIBackendConfig struct {
	// MountPath is not decoded directly by using the struct tags, and is instead populated by
	// ParseConfig when it is initialised
	MountPath string
	// Version is not decoded directly by using the struct tags, and is instead populated by ParseConfig
	// when it is initialised
	Version string

	Roles []Role `hcl:"role,block"`
}

type Role struct {
	Name      string                      `hcl:"role,label"`
	Zone      string                      `hcl:"zone,optional"`
	Secret    venafi.VenafiSecret         `hcl:"secret,block"`
	TestCerts []venafi.CertificateRequest `hcl:"test_certificate,block"`
}

func (c *VenafiPKIBackendConfig) ValidateConfig() error {
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
	err := r.Secret.Validate(venafi.SecretsEngine)
	if err != nil {
		return err
	}

	return nil
}

func (r *Role) WriteHCL(hclBody *hclwrite.Body) {
	roleBlock := hclBody.AppendNewBlock("role", []string{r.Name})
	roleBody := roleBlock.Body()
	r.Secret.WriteHCL(roleBody)
}

func (c *VenafiPKIBackendConfig) GenerateConfigAndWriteHCL(hclBody *hclwrite.Body) error {
	for i := 1; true; i++ {
		role, err := askForRole()
		if err != nil {
			return err
		}
		role.WriteHCL(hclBody)

		question := questions.ClosedQuestion{
			Question: fmt.Sprintf("You have configured %d roles, are there more", i),
			Items:    []string{"Yes", "No that's it"},
		}
		answer, err := question.Ask()
		if answer[0] != "Yes" {
			break
		}

		hclBody.AppendNewline()
	}
	// TODO: test certs (loop)
	return nil
}

func askForRole() (*Role, error) {
	answers, err := questions.AskQuestions([]questions.Question{
		&questions.OpenEndedQuestion{
			Question: "What should the role be called?",
		},
		&questions.QuestionBranch{
			ConditionQuestion: &questions.ClosedQuestion{
				Question: "What type of Venafi instance will be used?",
				Items:    []string{"TPP", "Venafi-as-a-Service"},
			},
			ConditionAnswer: "TPP",
			BranchA: []questions.Question{
				&questions.OpenEndedQuestion{
					Question: "What is the URL of the TPP instance?",
					Default:  "$TPP_URL",
				},
				&questions.OpenEndedQuestion{
					Question: "What is the username used to access the TPP instance?",
					Default:  "$TPP_USERNAME",
				},
				&questions.OpenEndedQuestion{
					Question: "What is the password of the TPP user?",
					Default:  "$TPP_PASSWORD",
				},
			},
			BranchB: []questions.Question{
				&questions.OpenEndedQuestion{
					Question: "What is the Venafi-as-a-Service API Key?",
					Default:  "$VENAFI_APIKEY",
				},
			},
		},
		&questions.OpenEndedQuestion{
			Question: "What project zone should be used for issuing certificates?",
		},
	})
	if err != nil {
		return nil, err
	}

	role := &Role{
		Name: string(answers[0]),
	}
	switch answers[1] {
	case "TPP":
		role.Secret.Name = "tpp"
		role.Secret.TPP = &venafi.VenafiTPPConnection{
			URL:      string(answers[2]),
			Username: string(answers[3]),
			Password: string(answers[4]),
			Zone:     string(answers[5]),
		}
	case "Venafi-as-a-Service":
		role.Secret.Name = "vaas"
		role.Secret.Cloud = &venafi.VenafiCloudConnection{
			APIKey: string(answers[2]),
			Zone:   string(answers[3]),
		}
	default:
		panic("unimplemented Venafi secret type, expected TPP or Venafi-as-a-Service")
	}
	return role, nil
}
