package kubernetes

import (
	"context"
	"encoding/base64"
	"encoding/json"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DockerCredential provides login information for authenticating to a Docker Registry
type DockerCredential struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Auth     string `json:"auth"`
}

// Endpoint provides information for authenticating to a specific Docker Registry Endpoint
type Endpoint map[string]*DockerCredential

// DockerSecret provides a structure for authentication information stored in a dockercfg (i.e. $HOME/.docker/config.json)
type DockerSecret struct {
	Auths *Endpoint `json:"auths"`
}

// NewDockerSecret returns an initialized Docker Registry secret
func NewDockerSecret(secretName, namespace, username, password, proxyEndpoint string) (*v1.Secret, error) {
	dockerSecret := DockerSecret{
		Auths: &Endpoint{
			proxyEndpoint: &DockerCredential{
				Username: username,
				Password: password,
				Auth:     base64.StdEncoding.EncodeToString([]byte(username + ":" + password)),
			},
		},
	}

	config, err := json.Marshal(dockerSecret)
	if err != nil {
		return nil, err
	}

	return &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespace,
		},
		Type: v1.SecretTypeDockerConfigJson,
		Data: map[string][]byte{
			v1.DockerConfigJsonKey: config,
		},
	}, nil
}

// UpdateSecret provides a method for updating a Kubernetes secret by either direct update or indirect update via a delete + create operation
func (c *Client) UpdateSecret(secret *v1.Secret) error {
	// Attempt update on secret
	_, err := c.CoreV1().Secrets(secret.ObjectMeta.Namespace).Update(context.TODO(), secret, metav1.UpdateOptions{})
	if err == nil {
		return nil
	}

	// Replace secret if update fails
	err = c.replaceSecret(secret)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) replaceSecret(secret *v1.Secret) (err error) {
	// Delete secret
	// No need to check for error here. We only care that the secret does not exist.
	c.CoreV1().Secrets(secret.ObjectMeta.Namespace).Delete(context.TODO(), secret.ObjectMeta.Name, metav1.DeleteOptions{})

	// Re-create secret
	_, err = c.CoreV1().Secrets(secret.ObjectMeta.Namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}
