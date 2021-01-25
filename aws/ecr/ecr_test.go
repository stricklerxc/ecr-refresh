package ecr

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
)

var response ecr.GetAuthorizationTokenOutput

// Mocks ECR API calls for GetAuthorizationToken
type mockGetAuthorizationToken struct {
	ecriface.ECRAPI
	Resp ecr.GetAuthorizationTokenOutput
}

func init() {
	// Setup mock response
	response = ecr.GetAuthorizationTokenOutput{
		AuthorizationData: []*ecr.AuthorizationData{
			{
				AuthorizationToken: aws.String("QVdTOnBhc3N3b3JkCg=="), // Base64 encoded "AWS:password"
				ProxyEndpoint:      aws.String("https://default.dkr.ecr.us-east-1.amazonaws.com"),
			},
		},
	}
}

func (m mockGetAuthorizationToken) GetAuthorizationToken(*ecr.GetAuthorizationTokenInput) (*ecr.GetAuthorizationTokenOutput, error) {
	return &m.Resp, nil
}

func Test_GetCredential_UpdatesReceiverWithCredential(t *testing.T) {
	// Create new registry
	registry := Registry{
		Client: mockGetAuthorizationToken{Resp: response},
	}

	// Test that Credentials is nil
	if registry.Credential != nil {
		t.Errorf("Expected credential to be nil. Got: %v", registry.Credential)
	}

	// Get credential
	err := registry.GetCredential()
	if err != nil {
		t.Fatalf("registry.GetCredential() call failed. Caused by: %v", err.Error())
	}

	// Test that registry credential was updated
	if registry.Credential.Username != "AWS" {
		t.Errorf("Expected username to be \"AWS\". Got: \"%v\"", registry.Credential.Username)
	}

	if registry.Credential.Password != "password" {
		t.Errorf("Expected password to be \"password\". Got: \"%v\"", registry.Credential.Password)
	}

	if registry.Credential.ProxyEndpoint != "https://default.dkr.ecr.us-east-1.amazonaws.com" {
		t.Errorf("Expected proxyendpoint to be \"https://default.dkr.ecr.us-east-1.amazonaws.com\". Got: \"%v\"", registry.Credential.ProxyEndpoint)
	}
}
