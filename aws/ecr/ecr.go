package ecr

import (
	"encoding/base64"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
)

// Registry provides the ECR API client to make ECR operational calls, as well as, credentials
// for authenticating to ECR Registries with IAM credentials
type Registry struct {
	Client     ecriface.ECRAPI
	Credential *Credential
}

// Credential provides authentication information for ECR Registries
type Credential struct {
	Username      string
	Password      string
	ProxyEndpoint string
}

// NewRegistry returns a Registry type initialized with an ECR API client
func NewRegistry(config *aws.Config) *Registry {
	return &Registry{
		Client:     ecr.New(session.Must(session.NewSession()), config),
		Credential: nil,
	}
}

// GetCredential updates a Registry type with credential information for authenticating to ECR Registries
func (r *Registry) GetCredential() error {
	resp, err := r.Client.GetAuthorizationToken(&ecr.GetAuthorizationTokenInput{})
	if err != nil {
		return err
	}

	tokenDecoded, err := base64.StdEncoding.DecodeString(*resp.AuthorizationData[0].AuthorizationToken)
	if err != nil {
		return err
	}

	userPass := strings.Split(strings.TrimSpace(string(tokenDecoded)), ":")

	r.Credential = &Credential{
		Username:      userPass[0],
		Password:      userPass[1],
		ProxyEndpoint: *resp.AuthorizationData[0].ProxyEndpoint,
	}

	return nil
}
