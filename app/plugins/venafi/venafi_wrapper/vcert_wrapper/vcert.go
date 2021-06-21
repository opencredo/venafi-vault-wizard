package vcert_wrapper

import (
	"fmt"
	"github.com/Venafi/vcert/v4"
	"github.com/Venafi/vcert/v4/pkg/certificate"
	"github.com/Venafi/vcert/v4/pkg/endpoint"
	"github.com/Venafi/vcert/v4/pkg/venafi/cloud"
	"github.com/Venafi/vcert/v4/pkg/venafi/tpp"
	"github.com/opencredo/venafi-vault-wizard/app/plugins/venafi"
	"github.com/opencredo/venafi-vault-wizard/app/plugins/venafi/venafi_wrapper"
	"net/http"
)

type venafiClient struct {
	cloudConnector *cloud.Connector
	tppConnector   *tpp.Connector
}

func NewVenafiClient(secret venafi.VenafiSecret, zone string) (venafi_wrapper.VenafiWrapper, error) {
	var config vcert.Config

	if secret.Cloud != nil {
		config = vcert.Config{
			ConnectorType: endpoint.ConnectorTypeCloud,
			Credentials: &endpoint.Authentication{
				APIKey: secret.Cloud.APIKey,
			},
			Zone: zone,
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
			Zone: zone,
		}
	}

	if config.ConnectorType == endpoint.ConnectorTypeCloud {
		client, err := vcert.NewClient(&config)
		if err != nil {
			return nil, err
		}
		cloudConnector, ok := client.(*cloud.Connector)
		if !ok {
			return nil, fmt.Errorf("cannot cast client to Cloud Connector")
		}

		return &venafiClient{
			cloudConnector: cloudConnector,
		}, nil
	} else if config.ConnectorType == endpoint.ConnectorTypeTPP {
		client, err := vcert.NewClient(&config)
		if err != nil {
			return nil, err
		}
		tppConnector, ok := client.(*tpp.Connector)
		if !ok {
			return nil, fmt.Errorf("cannot cast client to TPP Connector")
		}

		return &venafiClient{
			tppConnector: tppConnector,
		}, nil
	} else {
		panic("expected venafiClient to have either cloudConnector or tppConnector specified")
	}
}

func (v *venafiClient) setZone(zone string) {
	if v.cloudConnector != nil {
		v.cloudConnector.SetZone(zone)
	} else if v.tppConnector != nil {
		v.tppConnector.SetZone(zone)
	} else {
		panic("expected venafiClient to have either cloudConnector or tppConnector specified")
	}
}

func (v *venafiClient) GenerateRequest(config *endpoint.ZoneConfiguration, req *certificate.Request, zone string) (err error) {
	v.setZone(zone)
	if v.cloudConnector != nil {
		return v.cloudConnector.GenerateRequest(config, req)
	} else if v.tppConnector != nil {
		return v.tppConnector.GenerateRequest(config, req)
	} else {
		panic("expected venafiClient to have either cloudConnector or tppConnector specified")
	}
}

func (v *venafiClient) RequestCertificate(req *certificate.Request, zone string) (requestID string, err error) {
	v.setZone(zone)
	if v.cloudConnector != nil {
		return v.cloudConnector.RequestCertificate(req)
	} else if v.tppConnector != nil {
		return v.tppConnector.RequestCertificate(req)
	} else {
		panic("expected venafiClient to have either cloudConnector or tppConnector specified")
	}
}

func (v *venafiClient) RetrieveCertificate(req *certificate.Request, zone string) (certificates *certificate.PEMCollection, err error) {
	v.setZone(zone)
	if v.cloudConnector != nil {
		return v.cloudConnector.RetrieveCertificate(req)
	} else if v.tppConnector != nil {
		return v.tppConnector.RetrieveCertificate(req)
	} else {
		panic("expected venafiClient to have either cloudConnector or tppConnector specified")
	}
}

func (v *venafiClient) GetRefreshToken(auth *endpoint.Authentication) (resp tpp.OauthGetRefreshTokenResponse, err error) {
	if v.tppConnector != nil {
		return v.tppConnector.GetRefreshToken(auth)
	} else {
		panic("expected venafiClient to have either cloudConnector or tppConnector specified")
	}
}
