package test

import (
	"github.com/boeboe/iec61850"
	"testing"
	"time"
)

func CreateClient(t *testing.T) *iec61850.Client {
	settings := iec61850.NewSettings()

	client, err := iec61850.NewClient(settings)
	if err != nil {
		t.Fatalf("create client error %v\n", err)
	}
	return client
}

func CreateTlsClient(t *testing.T) *iec61850.Client {
	settings := iec61850.NewSettings()
	settings.Port = -1
	tlsConfig := iec61850.NewTLSConfig()

	tlsConfig.KeyFile = "client_CA1_1.key"
	tlsConfig.CertFile = "client_CA1_1.pem"
	tlsConfig.ChainValidation = true
	tlsConfig.AllowOnlyKnownCertificates = false
	tlsConfig.AddCACertificateFromFile("root_CA1.pem")

	client, err := iec61850.NewClientWithTlsSupport(settings, tlsConfig)
	if err != nil {
		t.Fatalf("create client error %v\n", err)
	}
	return client
}

func CloseClient(client *iec61850.Client) {
	time.Sleep(time.Second)
	client.Close()
}

func DoRead(t *testing.T, client *iec61850.Client, objectRef string, fc iec61850.FC) {
	value, err := client.Read(objectRef, fc)
	if err != nil {
		t.Fatalf("read %s object error %v\n", objectRef, err)
	}
	t.Logf("read %s value -> %v", objectRef, value)
}
