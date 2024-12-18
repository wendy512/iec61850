package server

import (
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/wendy512/iec61850"
)

func TestCreateServerFromConfigFile2(t *testing.T) {
	model, err := iec61850.CreateModelFromConfigFileEx("complexModel.cfg")
	if err != nil {
		t.Fatalf("create model error %v\n", err)
	}

	server := iec61850.NewServerWithConfig(iec61850.NewServerConfig(), model)

	modelNode := model.GetModelNodeByObjectReference("ied1Inverter/ZINV1.OutVarSet.setMag.f")
	server.SetHandleWriteAccess(modelNode, func(node *iec61850.ModelNode, mmsValue *iec61850.MmsValue) iec61850.MmsDataAccessError {
		t.Logf("handle write access, value %#v\n", mmsValue)
		return iec61850.DATA_ACCESS_ERROR_SUCCESS
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
