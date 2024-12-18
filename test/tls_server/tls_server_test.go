package tls_server

import (
	"github.com/wendy512/iec61850"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"unsafe"
)

func TestTlsServer(t *testing.T) {
	model, err := iec61850.CreateModelFromConfigFileEx("model.cfg")
	if err != nil {
		t.Fatalf("create model error %v\n", err)
	}

	tlsConfig := iec61850.NewTLSConfig()
	tlsConfig.KeyFile = "server_CA1_1.key"
	tlsConfig.CertFile = "server_CA1_1.pem"
	tlsConfig.AddCACertificateFromFile("root_CA1.pem")
	tlsConfig.AddAllowedCertificateFromFile("client_CA1_1.pem")
	tlsConfig.AddAllowedCertificateFromFile("client_CA1_2.pem")
	tlsConfig.AllowOnlyKnownCertificates = true
	tlsConfig.ChainValidation = false

	server, err := iec61850.NewServerWithTlsSupport(iec61850.NewServerConfig(), tlsConfig, model)

	modelNode := model.GetModelNodeByObjectReference("simpleIOGenericIO/GGIO1.SPCSO1")
	server.SetControlHandler(modelNode, func(node *iec61850.ModelNode, action *iec61850.ControlAction, mmsValue *iec61850.MmsValue, test bool) iec61850.ControlHandlerResult {
		t.Logf("control handler, action %#v, mmsValue %#v\n", action, mmsValue)
		return iec61850.CONTROL_RESULT_OK
	})
	server.SetAuthenticator(func(securityToken *unsafe.Pointer, authParameter *iec61850.AcseAuthenticationParameter, appReference *iec61850.IsoApplicationReference) bool {
		return true
	})

	server.Start(-1)
	defer server.Destroy()
	t.Logf("Server start up\n")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	<-sig
	t.Logf("Server stop\n")
	server.Stop()
}
