package venafi

import (
	"fmt"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/opencredo/venafi-vault-wizard/app/questions"
	"github.com/zclconf/go-cty/cty"
	"time"
)

type OptionalConfig struct {
	GenerateLease bool   `hcl:"generate_lease,optional"`
	AllowAnyName  bool   `hcl:"allow_any_name,optional"`
	TTL           string `hcl:"ttl,optional"`
	MaxTTL        string `hcl:"max_ttl,optional"`
}

func GenerateOptionalQuestions(questioner questions.Questioner) (*OptionalConfig, error) {
	q := map[string]questions.Question{
		"optional_parameters": questioner.NewClosedQuestion(&questions.ClosedQuestion{
			Question: "Do you want to configure optional parameters?",
			Items:    []string{"Yes", "No"},
		}),
		"generate_lease": questioner.NewClosedQuestion(&questions.ClosedQuestion{
			Question: "Do you want Vault to generate leases?",
			Items:    []string{"Yes", "No"},
		}),
		"allow_any_name": questioner.NewClosedQuestion(&questions.ClosedQuestion{
			Question: "Do you want Vault to allow any name?",
			Items:    []string{"Yes", "No"},
		}),
		"ttl": questioner.NewOpenEndedQuestion(&questions.OpenEndedQuestion{
			Question: "What should the default TTL be? (blank to use system default)",
			Default:  "",
		}),
		"max_ttl": questioner.NewOpenEndedQuestion(&questions.OpenEndedQuestion{
			Question: "What should the max TTL be? (blank to use system default)",
			Default:  "",
		}),
	}

	err := questions.AskQuestions([]questions.Question{
		&questions.QuestionBranch{
			ConditionQuestion: q["optional_parameters"],
			ConditionAnswer:   "Yes",
			BranchA: []questions.Question{
				q["generate_lease"],
				q["allow_any_name"],
				q["ttl"],
				q["max_ttl"],
			},
		},
	})
	if err != nil {
		return nil, err
	}

	if q["optional_parameters"].Answer() == "Yes" {
		return parseAnswers(q), nil
	}

	return nil, nil

}

func (oc *OptionalConfig) WriteHCL(hclBody *hclwrite.Body) {
	if oc.GenerateLease == true {
		hclBody.SetAttributeValue("generate_lease", cty.BoolVal(true))
	}

	if oc.AllowAnyName == true {
		hclBody.SetAttributeValue("allow_any_name", cty.BoolVal(true))
	}

	if oc.TTL != "" {
		hclBody.SetAttributeValue("ttl", cty.StringVal(oc.TTL))
	}

	if oc.MaxTTL != "" {
		hclBody.SetAttributeValue("max_ttl", cty.StringVal(oc.MaxTTL))
	}
}

func parseAnswers(q map[string]questions.Question) *OptionalConfig {
	optionalConfig := OptionalConfig{}

	if q["generate_lease"].Answer() == "Yes" {
		optionalConfig.GenerateLease = true
	}
	if q["allow_any_name"].Answer() == "Yes" {
		optionalConfig.AllowAnyName = true
	}

	optionalConfig.TTL = string(q["ttl"].Answer())
	optionalConfig.MaxTTL = string(q["max_ttl"].Answer())

	return &optionalConfig
}

func (oc *OptionalConfig) Validate() error {
	var TTL, MaxTTL time.Duration

	defaultDuration, err := time.ParseDuration("768h")
	if err != nil {
		return err
	}

	if oc.TTL == "" {
		TTL = defaultDuration
	} else {
		TTL, err = time.ParseDuration(oc.TTL)
		if err != nil {
			return fmt.Errorf("cannot parse ttl: %s", err)
		}
	}

	if oc.MaxTTL == "" {
		MaxTTL = defaultDuration
	} else {
		MaxTTL, err = time.ParseDuration(oc.MaxTTL)
		if err != nil {
			return fmt.Errorf("cannot parse max_ttl: %s", err)
		}
	}

	if MaxTTL < TTL {
		return fmt.Errorf("max_ttl must be greater than or equal to ttl")
	}

	return nil
}

func (oc *OptionalConfig) GetAsMap() map[string]interface{} {
	return map[string]interface{}{
		"ttl":            oc.TTL,
		"max_ttl":        oc.MaxTTL,
		"allow_any_name": oc.AllowAnyName,
		"generate_lease": oc.GenerateLease,
	}
}
