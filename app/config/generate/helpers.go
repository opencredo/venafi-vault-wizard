package generate

import (
	"strings"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

func WriteStringAttributeToHCL(attributeName, input string, body *hclwrite.Body) {
	if strings.HasPrefix(input, "$") {
		body.SetAttributeRaw(attributeName, envFunctionTokens(input[1:]))
		return
	}

	body.SetAttributeValue(attributeName, cty.StringVal(input))
}

func envFunctionTokens(environmentVariableName string) hclwrite.Tokens {
	return hclwrite.Tokens{
		&hclwrite.Token{
			Type:  hclsyntax.TokenIdent,
			Bytes: []byte("env"),
		},
		&hclwrite.Token{
			Type:  hclsyntax.TokenOParen,
			Bytes: []byte("("),
		},
		&hclwrite.Token{
			Type:  hclsyntax.TokenQuotedLit,
			Bytes: []byte("\"" + environmentVariableName + "\""),
		},
		&hclwrite.Token{
			Type:  hclsyntax.TokenCParen,
			Bytes: []byte(")"),
		},
	}
}
