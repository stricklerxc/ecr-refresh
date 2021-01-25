package kubernetes

import (
	"context"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

var fakeclient *Client

func init() {
	// Initialize fakeclient
	fakeclient = &Client{fake.NewSimpleClientset(&v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testSecret",
			Namespace: "default",
		},
		Data: map[string][]byte{
			".dockerconfigjson": []byte(""),
		},
	}), nil}
}

func Test_NewSecret_InitializesSecretSuccessfully(t *testing.T) {
	// Initialize secret
	secret, err := NewDockerSecret("testSecret", "default", "username", "password", "xxxxxxx.docker-registry.local")

	// Tests to ensure secret initialized properly
	if err != nil {
		t.Fatalf("NewDockerSecret() call failed. Caused by: %v", err.Error())
	}

	if secret.ObjectMeta.Name != "testSecret" {
		t.Errorf("Expected secret.ObjectMeta.Name to be \"testSecret\". Got: \"%v\"", secret.ObjectMeta.Name)
	}

	if secret.ObjectMeta.Namespace != "default" {
		t.Errorf("Expected secret.ObjectMeta.Namespace to be \"default\". Got: \"%v\"", secret.ObjectMeta.Namespace)
	}

	if string(secret.Data[".dockerconfigjson"]) == "" {
		t.Error("Expected secret.Object.Data[\".dockerconfigjson\"] to be populated but got nil.")
	}
}

func Test_Update_UpdateSecretSuccessfully(t *testing.T) {
	// Initialize new secret data
	secret, err := NewDockerSecret("testSecret", "default", "username", "password", "xxxxxxx.docker-registry.local")
	if err != nil {
		t.Fatalf("NewDockerSecret() call for secret failed. Caused by: %v", err.Error())
	}

	// Initialize secret
	oldSecret, err := fakeclient.CoreV1().Secrets("default").Get(context.TODO(), "testSecret", metav1.GetOptions{})

	// Update secret data
	err = fakeclient.UpdateSecret(secret)
	if err != nil {
		t.Fatalf("client.UpdateSecret() call failed. Caused by: %v", err.Error())
	}

	// Get updated secret
	newSecret, err := fakeclient.CoreV1().Secrets("default").Get(context.TODO(), "testSecret", metav1.GetOptions{})

	// Test secret contains new data
	if string(newSecret.Data[".dockerconfigjson"]) == "null" {
		t.Errorf("Expected newSecret.Data[\".dockerconfigjson\"] to be not \"null\". Got %v.", string(newSecret.Data[".dockerconfigjson"]))
	}

	if string(newSecret.Data[".dockerconfigjson"]) == string(oldSecret.Data[".dockerconfigjson"]) {
		t.Error("New secret should have different data than old secret.")
	}
}
