package kubernetes

import (
	"fmt"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// Client provides a Kubernetes Client & Config for interacting with the Kubernetes API
type Client struct {
	kubernetes.Interface
	Config *rest.Config
}

// NewClient returns an initialized client from either the configured service account or from $HOME/.kube/config
func NewClient() (*Client, error) {
	// Get config from /var/run/secrets/kubernetes.io/serviceaccount/token
	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Printf("INFO: %v", err.Error())

		// Get config from ~/.kube/config, if not in cluster
		config, err = clientcmd.BuildConfigFromFlags("", filepath.Join(homedir.HomeDir(), ".kube", "config"))
		if err != nil {
			return nil, err
		}
	}

	// Initialize client
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Client{
		client,
		config,
	}, nil
}
