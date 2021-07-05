package pki_backend

import (
	"fmt"
	"github.com/opencredo/venafi-vault-wizard/app/plugins"

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
	// BuildArch allows defining the build architecture
	BuildArch string `hcl:"build_arch,optional"`

	Roles []Role `hcl:"role,block"`
}

type Role struct {
	Name      string                      `hcl:"role,label"`
	Zone      string                      `hcl:"zone,optional"`
	Secret    venafi.VenafiSecret         `hcl:"secret,block"`
	TestCerts []venafi.CertificateRequest `hcl:"test_certificate,block"`

	OptionalConfig *venafi.OptionalConfig `hcl:"optional_config,block"`
}

func (c *VenafiPKIBackendConfig) ValidateConfig() error {
	err := plugins.ValidateBuildArch(c.BuildArch)
	if err != nil {
		return err
	}
	if len(c.Roles) == 0 {
		return fmt.Errorf("error at least one role must be provided: %w", errors.ErrBlankParam)
	}
	for _, role := range c.Roles {
		err = role.Validate()
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

	if r.OptionalConfig != nil {
		err = r.OptionalConfig.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Role) WriteHCL(hclBody *hclwrite.Body) {
	roleBlock := hclBody.AppendNewBlock("role", []string{r.Name})
	roleBody := roleBlock.Body()
	r.Secret.WriteHCL(roleBody)

	if r.OptionalConfig != nil {
		roleBody.AppendNewline()
		optionalBlock := roleBody.AppendNewBlock("optional_config", nil)
		r.OptionalConfig.WriteHCL(optionalBlock.Body())
	}

	for _, testCert := range r.TestCerts {
		roleBody.AppendNewline()
		certBlock := roleBody.AppendNewBlock("test_certificate", nil)
		testCert.WriteHCL(certBlock.Body())
	}
}

func (c *VenafiPKIBackendConfig) GenerateConfigAndWriteHCL(questioner questions.Questioner, hclBody *hclwrite.Body) error {
	for i := 1; true; i++ {
		role, err := askForRole(questioner)
		if err != nil {
			return err
		}

		hclBody.AppendNewline()
		role.WriteHCL(hclBody)

		moreRolesQuestion := questioner.NewClosedQuestion(&questions.ClosedQuestion{
			Question: fmt.Sprintf("You have configured %d roles, are there more", i),
			Items:    []string{"Yes", "No that's it"},
		})
		err = moreRolesQuestion.Ask()
		if err != nil {
			return err
		}
		if moreRolesQuestion.Answer() != "Yes" {
			break
		}
	}
	return nil
}

func askForRole(questioner questions.Questioner) (*Role, error) {
	q := map[string]questions.Question{
		"role_name": questioner.NewOpenEndedQuestion(&questions.OpenEndedQuestion{
			Question: "What should the role be called?",
		}),
		"venafi_type": questioner.NewClosedQuestion(&questions.ClosedQuestion{
			Question: "What type of Venafi instance will be used?",
			Items:    []string{"TPP", "Venafi-as-a-Service"},
		}),
		"tpp_url": questioner.NewOpenEndedQuestion(&questions.OpenEndedQuestion{
			Question: "What is the URL of the TPP instance?",
			Default:  "$TPP_URL",
		}),
		"tpp_user": questioner.NewOpenEndedQuestion(&questions.OpenEndedQuestion{
			Question: "What is the username used to access the TPP instance?",
			Default:  "$TPP_USERNAME",
		}),
		"tpp_password": questioner.NewOpenEndedQuestion(&questions.OpenEndedQuestion{
			Question: "What is the password of the TPP user?",
			Default:  "$TPP_PASSWORD",
		}),
		"apikey": questioner.NewOpenEndedQuestion(&questions.OpenEndedQuestion{
			Question: "What is the Venafi-as-a-Service API Key?",
			Default:  "$VENAFI_APIKEY",
		}),
		"zone": questioner.NewOpenEndedQuestion(&questions.OpenEndedQuestion{
			Question: "What project zone should be used for issuing certificates?",
		}),
	}
	err := questions.AskQuestions([]questions.Question{
		q["role_name"],
		&questions.QuestionBranch{
			ConditionQuestion: q["venafi_type"],
			ConditionAnswer:   "TPP",
			BranchA: []questions.Question{
				q["tpp_url"],
				q["tpp_user"],
				q["tpp_password"],
			},
			BranchB: []questions.Question{
				q["apikey"],
			},
		},
		q["zone"],
	})
	if err != nil {
		return nil, err
	}

	role := &Role{
		Name: string(q["role_name"].Answer()),
	}
	switch q["venafi_type"].Answer() {
	case "TPP":
		role.Secret.Name = "tpp"
		role.Secret.TPP = &venafi.VenafiTPPConnection{
			URL:      string(q["tpp_url"].Answer()),
			Username: string(q["tpp_user"].Answer()),
			Password: string(q["tpp_password"].Answer()),
			Zone:     string(q["zone"].Answer()),
		}
	case "Venafi-as-a-Service":
		role.Secret.Name = "vaas"
		role.Secret.VAAS = &venafi.VenafiVAASConnection{
			APIKey: string(q["apikey"].Answer()),
			Zone:   string(q["zone"].Answer()),
		}
	default:
		panic("unimplemented Venafi secret type, expected TPP or Venafi-as-a-Service")
	}

	optionalConfig, err := venafi.GenerateOptionalQuestions(questioner)
	if err != nil {
		return nil, err
	}
	role.OptionalConfig = optionalConfig

	testCertificates, err := venafi.AskForTestCertificates(questioner)
	if err != nil {
		return nil, err
	}

	role.TestCerts = testCertificates

	return role, nil
}
