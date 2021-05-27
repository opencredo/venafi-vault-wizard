package venafi

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/opencredo/venafi-vault-wizard/app/questions"
	"github.com/zclconf/go-cty/cty"
)

type CertificateRequest struct {
	CommonName   string `hcl:"common_name"`
	OU           string `hcl:"ou"`
	Organisation string `hcl:"organisation"`
	Locality     string `hcl:"locality"`
	Province     string `hcl:"province"`
	Country      string `hcl:"country"`
	TTL          string `hcl:"ttl"`
}

func (c *CertificateRequest) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"common_name":  c.CommonName,
		"ou":           c.OU,
		"organization": c.Organisation,
		"locality":     c.Locality,
		"province":     c.Province,
		"country":      c.Country,
		"ttl":          c.TTL,
	}
}

func GenerateCertRequestConfig(questioner questions.Questioner) (*CertificateRequest, error) {
	q := map[string]questions.Question{
		"cn": questioner.NewOpenEndedQuestion(&questions.OpenEndedQuestion{
			Question: "What should the common name (CN) of the certificate be?",
		}),
		"ou": questioner.NewOpenEndedQuestion(&questions.OpenEndedQuestion{
			Question: "What should the organisational unit (OU) of the certificate be?",
		}),
		"o": questioner.NewOpenEndedQuestion(&questions.OpenEndedQuestion{
			Question: "What should the organisation (O) of the certificate be?",
		}),
		"l": questioner.NewOpenEndedQuestion(&questions.OpenEndedQuestion{
			Question: "What should the locality (L) of the certificate be?",
		}),
		"p": questioner.NewOpenEndedQuestion(&questions.OpenEndedQuestion{
			Question: "What should the province (P) of the certificate be?",
		}),
		"c": questioner.NewOpenEndedQuestion(&questions.OpenEndedQuestion{
			Question: "What should the country (C) of the certificate be?",
		}),
		"ttl": questioner.NewOpenEndedQuestion(&questions.OpenEndedQuestion{
			Question: "What should the time-to-live (TTL) of the certificate be?",
		}),
	}
	err := questions.AskQuestions([]questions.Question{
		q["cn"],
		q["ou"],
		q["o"],
		q["l"],
		q["p"],
		q["c"],
		q["ttl"],
	})
	if err != nil {
		return nil, err
	}

	return &CertificateRequest{
		CommonName:   string(q["cn"].Answer()),
		OU:           string(q["ou"].Answer()),
		Organisation: string(q["o"].Answer()),
		Locality:     string(q["l"].Answer()),
		Province:     string(q["p"].Answer()),
		Country:      string(q["c"].Answer()),
		TTL:          string(q["ttl"].Answer()),
	}, nil
}

func (c *CertificateRequest) WriteHCL(hclBody *hclwrite.Body) {
	hclBody.SetAttributeValue("common_name", cty.StringVal(c.CommonName))
	hclBody.SetAttributeValue("ou", cty.StringVal(c.OU))
	hclBody.SetAttributeValue("organisation", cty.StringVal(c.Organisation))
	hclBody.SetAttributeValue("locality", cty.StringVal(c.Locality))
	hclBody.SetAttributeValue("province", cty.StringVal(c.Province))
	hclBody.SetAttributeValue("country", cty.StringVal(c.Country))
	hclBody.SetAttributeValue("ttl", cty.StringVal(c.TTL))
}

func AskForTestCertificates(questioner questions.Questioner) ([]CertificateRequest, error) {
	anyCSRsQuestion := questioner.NewClosedQuestion(&questions.ClosedQuestion{
		Question: "Would you like to request any test certificates to check everything is working?",
		Items:    []string{"Yes", "No, skip"},
	})
	err := anyCSRsQuestion.Ask()
	if err != nil {
		return nil, err
	}

	if anyCSRsQuestion.Answer() == "No, skip" {
		return nil, nil
	}

	var csrs []CertificateRequest
	for i := 1; true; i++ {
		certificateRequest, err := GenerateCertRequestConfig(questioner)
		if err != nil {
			return nil, err
		}

		csrs = append(csrs, *certificateRequest)

		moreCSRsQuestion := questioner.NewClosedQuestion(&questions.ClosedQuestion{
			Question: fmt.Sprintf("You have configured %d test certificates, are there more", i),
			Items:    []string{"Yes", "No that's it"},
		})
		err = moreCSRsQuestion.Ask()
		if err != nil {
			return nil, err
		}

		if moreCSRsQuestion.Answer() != "Yes" {
			break
		}
	}

	return csrs, nil
}
