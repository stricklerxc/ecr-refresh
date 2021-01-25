package main

import (
	"github.com/aws/aws-sdk-go/aws"

	"github.com/stricklerxc/ecr-refresh/aws/ecr"
	"github.com/stricklerxc/ecr-refresh/kubernetes"
)

func main() {
	// Initialize ECR Registry
	registry := ecr.NewRegistry(aws.NewConfig())

	// Initialize K8s Client
	client, err := kubernetes.NewClient()
	if err != nil {
		panic(err.Error())
	}

	// Get ECR Credential
	if err = registry.GetCredential(); err != nil {
		panic(err.Error())
	}

	// Create Docker Registry secret
	secret, err := kubernetes.NewDockerSecret("ecr-creds", "default", registry.Credential.Username, registry.Credential.Password, registry.Credential.ProxyEndpoint)
	if err != nil {
		panic(err.Error())
	}

	// Update Secret
	if err = client.UpdateSecret(secret); err != nil {
		panic(err.Error())
	}
}
