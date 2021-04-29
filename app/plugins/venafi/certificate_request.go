package venafi

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
