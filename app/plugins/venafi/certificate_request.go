package venafi

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
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

func (c *CertificateRequest) WriteHCL(hclBody *hclwrite.Body) {
	hclBody.SetAttributeValue("common_name", cty.StringVal(c.CommonName))
	hclBody.SetAttributeValue("ou", cty.StringVal(c.OU))
	hclBody.SetAttributeValue("organisation", cty.StringVal(c.Organisation))
	hclBody.SetAttributeValue("locality", cty.StringVal(c.Locality))
	hclBody.SetAttributeValue("province", cty.StringVal(c.Province))
	hclBody.SetAttributeValue("country", cty.StringVal(c.Country))
	hclBody.SetAttributeValue("ttl", cty.StringVal(c.TTL))
}
