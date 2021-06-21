package venafi_wrapper

import (
	"github.com/Venafi/vcert/v4/pkg/certificate"
	"github.com/Venafi/vcert/v4/pkg/endpoint"
	"github.com/Venafi/vcert/v4/pkg/venafi/tpp"
)

type VenafiWrapper interface {
	GenerateRequest(config *endpoint.ZoneConfiguration, req *certificate.Request, zone string) (err error)
	RequestCertificate(req *certificate.Request, zone string) (requestID string, err error)
	RetrieveCertificate(req *certificate.Request, zone string) (certificates *certificate.PEMCollection, err error)
	// GetRefreshToken TPP implementation only
	GetRefreshToken(auth *endpoint.Authentication) (resp tpp.OauthGetRefreshTokenResponse, err error)
}
