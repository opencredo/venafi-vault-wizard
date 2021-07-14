package vcert_wrapper

import (
	"fmt"
	"net/http"

	"github.com/Venafi/vcert/v4"
	"github.com/Venafi/vcert/v4/pkg/certificate"
	"github.com/Venafi/vcert/v4/pkg/endpoint"
	"github.com/Venafi/vcert/v4/pkg/venafi/cloud"
	"github.com/Venafi/vcert/v4/pkg/venafi/tpp"
	"github.com/opencredo/venafi-vault-wizard/app/plugins/venafi"
	"github.com/opencredo/venafi-vault-wizard/app/plugins/venafi/venafi_wrapper"
)

type venafiClient struct {
	vaasConnector *cloud.Connector
	tppConnector  *tpp.Connector
}

func NewVenafiClient(secret venafi.VenafiSecret) (venafi_wrapper.VenafiWrapper, error) {
	var config vcert.Config

	if secret.VaaS != nil {
		config = vcert.Config{
			ConnectorType: endpoint.ConnectorTypeCloud,
			Credentials: &endpoint.Authentication{
				APIKey: secret.VaaS.APIKey,
			},
			Zone: "",
			// Specify the DefaultClient otherwise vcert_wrapper creates its own HTTP Client and for some reason this replaces
			// the TLSClientConfig with a non-nil value it gets from somewhere and breaks things with the following:
			// vcert_wrapper error: server error: server unavailable: Get "https://api.venafi.cloud/v1/useraccounts": net/http: HTTP/1.x transport connection broken: malformed HTTP response
			Client: http.DefaultClient,
		}
	} else {
		config = vcert.Config{
			ConnectorType: endpoint.ConnectorTypeTPP,
			BaseUrl:       secret.TPP.URL,
			Credentials: &endpoint.Authentication{
				User:     secret.TPP.Username,
				Password: secret.TPP.Password,
			},
			Zone: "",
		}
	}

	client, err := vcert.NewClient(&config)
	if err != nil {
		return nil, err
	}

	if config.ConnectorType == endpoint.ConnectorTypeCloud {
		vaasConnector, ok := client.(*cloud.Connector)
		if !ok {
			return nil, fmt.Errorf("cannot cast client to VaaS Connector")
		}

		return &venafiClient{
			vaasConnector: vaasConnector,
		}, nil
	} else if config.ConnectorType == endpoint.ConnectorTypeTPP {
		tppConnector, ok := client.(*tpp.Connector)
		if !ok {
			return nil, fmt.Errorf("cannot cast client to TPP Connector")
		}

		return &venafiClient{
			tppConnector: tppConnector,
		}, nil
	} else {
		panic("expected venafiClient to have either vaasConnector or tppConnector specified")
	}
}

func (v *venafiClient) setZone(zone string) {
	if v.vaasConnector != nil {
		v.vaasConnector.SetZone(zone)
	} else if v.tppConnector != nil {
		v.tppConnector.SetZone(zone)
	} else {
		panic("expected venafiClient to have either vaasConnector or tppConnector specified")
	}
}

func (v *venafiClient) GenerateRequest(config *endpoint.ZoneConfiguration, req *certificate.Request, zone string) (err error) {
	v.setZone(zone)
	if v.vaasConnector != nil {
		return v.vaasConnector.GenerateRequest(config, req)
	} else if v.tppConnector != nil {
		return v.tppConnector.GenerateRequest(config, req)
	} else {
		panic("expected venafiClient to have either vaasConnector or tppConnector specified")
	}
}

func (v *venafiClient) RequestCertificate(req *certificate.Request, zone string) (requestID string, err error) {
	v.setZone(zone)
	if v.vaasConnector != nil {
		return v.vaasConnector.RequestCertificate(req)
	} else if v.tppConnector != nil {
		return v.tppConnector.RequestCertificate(req)
	} else {
		panic("expected venafiClient to have either vaasConnector or tppConnector specified")
	}
}

func (v *venafiClient) RetrieveCertificate(req *certificate.Request, zone string) (certificates *certificate.PEMCollection, err error) {
	v.setZone(zone)
	if v.vaasConnector != nil {
		return v.vaasConnector.RetrieveCertificate(req)
	} else if v.tppConnector != nil {
		return v.tppConnector.RetrieveCertificate(req)
	} else {
		panic("expected venafiClient to have either vaasConnector or tppConnector specified")
	}
}

func (v *venafiClient) GetRefreshToken(auth *endpoint.Authentication) (resp tpp.OauthGetRefreshTokenResponse, err error) {
	if v.tppConnector != nil {
		return v.tppConnector.GetRefreshToken(auth)
	} else {
		panic("expected venafiClient to have either vaasConnector or tppConnector specified")
	}
}
