package pki_monitor

import (
	"fmt"

	"github.com/opencredo/venafi-vault-wizard/app/plugins"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/opencredo/venafi-vault-wizard/app/config/errors"
	"github.com/opencredo/venafi-vault-wizard/app/plugins/venafi"
	"github.com/opencredo/venafi-vault-wizard/app/questions"
	"github.com/zclconf/go-cty/cty"
)

type VenafiPKIMonitorConfig struct {
	// MountPath is not decoded directly by using the struct tags, and is instead populated by plugins.LookupPlugin
	// when it is initialised
	MountPath string
	// Version is not decoded directly by using the struct tags, and is instead populated by plugins.LookupPlugin
	// when it is initialised
	Version string
	// BuildArch allows defining the build architecture
	BuildArch string

	Role Role `hcl:"role,block"`
}

type Role struct {
	Name string `hcl:"role,label"`

	Secret UnZonedSecret `hcl:"secret,block"`

	EnforcementPolicy *Policy `hcl:"enforcement_policy,block"`
	ImportPolicy      *Policy `hcl:"import_policy,block"`

	IntermediateCert *IntermediateCertRequest   `hcl:"intermediate_certificate,block"`
	RootCert         *venafi.CertificateRequest `hcl:"root_certificate,block"`

	TestCerts []venafi.CertificateRequest `hcl:"test_certificate,block"`

	OptionalConfig *venafi.OptionalConfig `hcl:"optional_config,block"`
}

type IntermediateCertRequest struct {
	Zone                      string          `hcl:"zone"`
	venafi.CertificateRequest `hcl:",remain"` // gohcl currently ignores any field without hcl tags, even in an embedded struct with nested tagged fields
}

type Policy struct {
	Zone string `hcl:"zone"`
}

// UnZonedSecret Used to add the label, and to maintain consistent structure with other uses of VenafiSecret.
type UnZonedSecret struct {
	Name                string `hcl:"name,label"`
	venafi.VenafiSecret `hcl:",remain"`
}

func (c *VenafiPKIMonitorConfig) ValidateConfig() error {
	err := plugins.ValidateBuildArch(c.BuildArch)
	if err != nil {
		return err
	}
	err = c.Role.Validate()
	if err != nil {
		return err
	}
	return nil
}

func (r *Role) Validate() error {
	err := r.Secret.Validate()
	if err != nil {
		return err
	}

	if r.OptionalConfig != nil {
		err = r.OptionalConfig.Validate()
		if err != nil {
			return err
		}
	}

	intermediateCertProvided := r.IntermediateCert != nil
	rootCertProvided := r.RootCert != nil

	if (intermediateCertProvided && rootCertProvided) || (!intermediateCertProvided && !rootCertProvided) {
		return fmt.Errorf("error, must provide exactly one of either the intermediate_certificate or root_certificate blocks: %w", errors.ErrConflictingBlocks)
	}

	if r.EnforcementPolicy == nil && r.ImportPolicy == nil {
		return fmt.Errorf("error, at least one of either enforcement_policy or import_policy must be provided: %w", errors.ErrBlankParam)
	}

	if intermediateCertProvided && r.IntermediateCert.Zone == "" {
		return fmt.Errorf("error, intermediate_certificate zone cannot be an empty string: %w", errors.ErrBlankParam)
	}

	if r.EnforcementPolicy != nil && r.EnforcementPolicy.Zone == "" {
		return fmt.Errorf("error, enforcement_policy zone cannot be an empty string: %w", errors.ErrBlankParam)
	}

	if r.ImportPolicy != nil && r.ImportPolicy.Zone == "" {
		return fmt.Errorf("error, import_policy zone cannot be an empty string: %w", errors.ErrBlankParam)
	}

	return nil
}

func (r *Role) WriteHCL(hclBody *hclwrite.Body) {
	roleBlock := hclBody.AppendNewBlock("role", []string{r.Name})
	roleBody := roleBlock.Body()
	r.Secret.WriteHCL(roleBody)

	if r.EnforcementPolicy != nil {
		roleBody.AppendNewline()
		policyBlock := roleBody.AppendNewBlock("enforcement_policy", nil)
		policyBlock.Body().SetAttributeValue("zone", cty.StringVal(r.EnforcementPolicy.Zone))
	}

	if r.ImportPolicy != nil {
		roleBody.AppendNewline()
		policyBlock := roleBody.AppendNewBlock("import_policy", nil)
		policyBlock.Body().SetAttributeValue("zone", cty.StringVal(r.ImportPolicy.Zone))
	}

	if r.IntermediateCert != nil {
		roleBody.AppendNewline()
		certBlock := roleBody.AppendNewBlock("intermediate_certificate", nil)
		certBody := certBlock.Body()
		certBody.SetAttributeValue("zone", cty.StringVal(r.IntermediateCert.Zone))
		r.IntermediateCert.WriteHCL(certBody)
	}

	if r.RootCert != nil {
		roleBody.AppendNewline()
		certBlock := roleBody.AppendNewBlock("root_certificate", nil)
		r.RootCert.WriteHCL(certBlock.Body())
	}

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

func (s *UnZonedSecret) WriteHCL(hclBody *hclwrite.Body) {
	secretBlock := hclBody.AppendNewBlock("secret", []string{s.Name})
	secretBody := secretBlock.Body()
	s.VenafiSecret.WriteHCL(secretBody)
}

func (c *VenafiPKIMonitorConfig) GenerateConfigAndWriteHCL(questioner questions.Questioner, hclBody *hclwrite.Body) error {
	role, err := askForRole(questioner)
	if err != nil {
		return err
	}

	hclBody.AppendNewline()
	role.WriteHCL(hclBody)

	return nil
}

func askForRole(questioner questions.Questioner) (*Role, error) {
	q := map[string]questions.Question{
		"role_name": questioner.NewOpenEndedQuestion(&questions.OpenEndedQuestion{
			Question: "What should the role be called?",
		}),
		"venafi_type": questioner.NewClosedQuestion(&questions.ClosedQuestion{
			Question: "What type of Venafi instance will be used?",
			Items:    []string{"TPP", "Venafi as a Service"},
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
			Question: "What is the Venafi as a Service API Key?",
			Default:  "$VENAFI_APIKEY",
		}),
		"enforce_policy": questioner.NewClosedQuestion(&questions.ClosedQuestion{
			Question: "Would you like Vault to enforce a certificate policy from Venafi?",
			Items:    []string{"Yes", "No, allow Vault to issue any certificate"},
		}),
		"enforcement_policy_name": questioner.NewOpenEndedQuestion(&questions.OpenEndedQuestion{
			Question: "Which Venafi policy should be used?",
		}),
		"reuse_policy": questioner.NewClosedQuestion(&questions.ClosedQuestion{
			Question: "Should the same policy be used to import policies for visibility?",
			Items:    []string{"Yes", "No, use a separate policy"},
		}),
		"import_policy_name": questioner.NewOpenEndedQuestion(&questions.OpenEndedQuestion{
			Question: "Which Venafi policy should be used for importing certificates into?",
		}),
		"issuing_cert_type": questioner.NewClosedQuestion(&questions.ClosedQuestion{
			Question: "What type of certificate should Vault use to issue certificates?",
			Items:    []string{"Self-signed root certificate", "Intermediate certificate issued by Venafi"},
		}),
		"subca_policy": questioner.NewOpenEndedQuestion(&questions.OpenEndedQuestion{
			Question: "Which policy should be used to issue the subordinate CA certificate?",
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
		&questions.QuestionBranch{
			ConditionQuestion: q["enforce_policy"],
			ConditionAnswer:   "Yes",
			BranchA: []questions.Question{
				q["enforcement_policy_name"],
				&questions.QuestionBranch{
					ConditionQuestion: q["reuse_policy"],
					ConditionAnswer:   "No, use a separate policy",
					BranchA: []questions.Question{
						q["import_policy_name"],
					},
				},
			},
			BranchB: []questions.Question{
				q["import_policy_name"],
			},
		},
		&questions.QuestionBranch{
			ConditionQuestion: q["issuing_cert_type"],
			ConditionAnswer:   "Self-signed root certificate",
			BranchA:           nil,
			BranchB: []questions.Question{
				q["subca_policy"],
			},
		},
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
		}
	case "Venafi as a Service":
		role.Secret.Name = "vaas"
		role.Secret.VaaS = &venafi.VenafiVaaSConnection{
			APIKey: string(q["apikey"].Answer()),
		}
	default:
		panic("unimplemented Venafi secret type, expected TPP or Venafi as a Service")
	}

	if q["enforce_policy"].Answer() == "Yes" {
		role.EnforcementPolicy = &Policy{
			Zone: string(q["enforcement_policy_name"].Answer()),
		}

		if q["reuse_policy"].Answer() != "Yes" {
			role.ImportPolicy = &Policy{
				Zone: string(q["import_policy_name"].Answer()),
			}
		}
	} else {
		role.ImportPolicy = &Policy{
			Zone: string(q["import_policy_name"].Answer()),
		}
	}

	if q["issuing_cert_type"].Answer() == "Self-signed root certificate" {
		certificateRequest, err := venafi.GenerateCertRequestConfig(questioner)
		if err != nil {
			return nil, err
		}
		role.RootCert = certificateRequest
	} else {
		certificateRequest, err := venafi.GenerateCertRequestConfig(questioner)
		if err != nil {
			return nil, err
		}
		role.IntermediateCert = &IntermediateCertRequest{
			Zone:               string(q["subca_policy"].Answer()),
			CertificateRequest: *certificateRequest,
		}
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
