package server

import (
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/wendy512/iec61850"
)

func TestCreateServerFromConfigFile1(t *testing.T) {
	model, err := iec61850.CreateModelFromConfigFileEx("simpleIO_control_tests.cfg")
	if err != nil {
		t.Fatalf("create model error %v\n", err)
	}

	server := iec61850.NewServerWithConfig(iec61850.NewServerConfig(), model)

	modelNode := model.GetModelNodeByObjectReference("simpleIOGenericIO/GGIO1.SPCSO1")
	server.SetControlHandler(modelNode, func(node *iec61850.ModelNode, action *iec61850.ControlAction, mmsValue *iec61850.MmsValue, test bool) iec61850.ControlHandlerResult {
		t.Logf("control handler, action %#v, mmsValue %#v\n", action, mmsValue)
		return iec61850.CONTROL_RESULT_OK
	})

	server.Start(102)
	defer server.Destroy()
	t.Logf("Server start up\n")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	<-sig
	t.Logf("Server stop\n")
	server.Stop()
}
