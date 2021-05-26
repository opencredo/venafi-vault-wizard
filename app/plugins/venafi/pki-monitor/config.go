package pki_monitor

import (
	"fmt"

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

	Role Role `hcl:"role,block"`
}

type Role struct {
	Name string `hcl:"role,label"`

	Secret venafi.VenafiSecret `hcl:"secret,block"`

	EnforcementPolicy *Policy `hcl:"enforcement_policy,block"`
	ImportPolicy      *Policy `hcl:"import_policy,block"`

	IntermediateCert *venafi.CertificateRequest `hcl:"intermediate_certificate,block"`
	RootCert         *venafi.CertificateRequest `hcl:"root_certificate,block"`

	TestCerts []venafi.CertificateRequest `hcl:"test_certificate,block"`

	GenerateLease bool   `hcl:"generate_lease,optional"`
	AllowAnyName  bool   `hcl:"allow_any_name,optional"`
	TTL           string `hcl:"ttl,optional"`
	MaxTTL        string `hcl:"max_ttl,optional"`
}

type Policy struct {
	Zone string `hcl:"zone"`
}

func (c *VenafiPKIMonitorConfig) ValidateConfig() error {
	err := c.Role.Validate()
	if err != nil {
		return err
	}
	return nil
}

func (r *Role) Validate() error {
	err := r.Secret.Validate(venafi.MonitorEngine)
	if err != nil {
		return err
	}

	if r.MaxTTL < r.TTL {
		return fmt.Errorf("max_ttl must be greater than or equal to ttl")
	}

	intermediateCertProvided := r.IntermediateCert != nil
	rootCertProvided := r.RootCert != nil

	if (intermediateCertProvided && rootCertProvided) || (!intermediateCertProvided && !rootCertProvided) {
		return fmt.Errorf("error, must provide exactly one of either the intermediate_certificate or root_certificate blocks: %w", errors.ErrConflictingBlocks)
	}

	if r.EnforcementPolicy == nil && r.ImportPolicy == nil {
		return fmt.Errorf("error, at least one of either enforcement_policy or import_policy must be provided: %w", errors.ErrBlankParam)
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
		r.IntermediateCert.WriteHCL(certBlock.Body())
	}

	if r.RootCert != nil {
		roleBody.AppendNewline()
		certBlock := roleBody.AppendNewBlock("root_certificate", nil)
		r.RootCert.WriteHCL(certBlock.Body())
	}
}

func (c *VenafiPKIMonitorConfig) GenerateConfigAndWriteHCL(questioner questions.Questioner, hclBody *hclwrite.Body) error {
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
	// TODO: test certs (loop)
	return nil
}

func askForRole(questioner questions.Questioner) (*Role, error) {
	q := map[string]questions.Question{
		"role_name": questioner.NewOpenEndedQuestion(&questions.OpenEndedQuestion{
			Question: "What should the role be called?",
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
		"cn": questioner.NewOpenEndedQuestion(&questions.OpenEndedQuestion{
			Question: "What should the common name (CN) of the issuing certificate be?",
		}),
		"ou": questioner.NewOpenEndedQuestion(&questions.OpenEndedQuestion{
			Question: "What should the organisational unit (OU) of the issuing certificate be?",
		}),
		"o": questioner.NewOpenEndedQuestion(&questions.OpenEndedQuestion{
			Question: "What should the organisation (O) of the issuing certificate be?",
		}),
		"l": questioner.NewOpenEndedQuestion(&questions.OpenEndedQuestion{
			Question: "What should the locality (L) of the issuing certificate be?",
		}),
		"p": questioner.NewOpenEndedQuestion(&questions.OpenEndedQuestion{
			Question: "Qhat should the province (P) of the issuing certificate be?",
		}),
		"c": questioner.NewOpenEndedQuestion(&questions.OpenEndedQuestion{
			Question: "What should the country (C) of the issuing certificate be?",
		}),
		"ttl": questioner.NewOpenEndedQuestion(&questions.OpenEndedQuestion{
			Question: "What should the time-to-live (TTL) of the issuing certificate be?",
		}),
	}
	err := questions.AskQuestions([]questions.Question{
		q["role_name"],
		q["tpp_url"],
		q["tpp_user"],
		q["tpp_password"],
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
		q["cn"],
		q["ou"],
		q["o"],
		q["l"],
		q["p"],
		q["c"],
		q["ttl"],

		// TODO: extra questions around default TTL, max TTL etc
		// TODO: test certs
	})
	if err != nil {
		return nil, err
	}

	role := &Role{
		Name: string(q["role_name"].Answer()),
		Secret: venafi.VenafiSecret{
			Name: "tpp",
			TPP: &venafi.VenafiTPPConnection{
				URL:      string(q["tpp_url"].Answer()),
				Username: string(q["tpp_user"].Answer()),
				Password: string(q["tpp_password"].Answer()),
			},
		},
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
		role.RootCert = answersToCSR(q)
	} else {
		_ = q["subca_policy"].Answer() // TODO: issuing cert policy zone
		role.IntermediateCert = answersToCSR(q)
	}

	return role, nil
}

func answersToCSR(questions map[string]questions.Question) *venafi.CertificateRequest {
	return &venafi.CertificateRequest{
		CommonName:   string(questions["cn"].Answer()),
		OU:           string(questions["ou"].Answer()),
		Organisation: string(questions["o"].Answer()),
		Locality:     string(questions["l"].Answer()),
		Province:     string(questions["p"].Answer()),
		Country:      string(questions["c"].Answer()),
		TTL:          string(questions["ttl"].Answer()),
	}
}
