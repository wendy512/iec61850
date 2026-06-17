package server

import (
	"testing"
	"time"

	"github.com/wendy512/iec61850"
)

const (
	pcControlRef = "simpleIOGenericIO/GGIO1.SPCSO1"
	pcStValRef   = pcControlRef + ".stVal"
)

func operateWithCheck(t *testing.T, port int, check iec61850.PerformCheckHandler) (error, bool) {
	t.Helper()

	model, err := iec61850.CreateModelFromConfigFileEx("simpleIO_control_tests.cfg")
	if err != nil {
		t.Fatalf("create model: %v", err)
	}
	server := iec61850.NewServerWithConfig(iec61850.NewServerConfig(), model)
	node := model.GetModelNodeByObjectReference(pcControlRef)

	server.SetPerformCheckHandler(node, check)
	server.SetControlHandler(node, func(_ *iec61850.ModelNode, _ *iec61850.ControlAction, _ *iec61850.MmsValue, _ bool) iec61850.ControlHandlerResult {
		return iec61850.CONTROL_RESULT_OK
	})

	server.Start(port)
	defer server.Stop()
	time.Sleep(300 * time.Millisecond)

	settings := iec61850.NewSettings()
	settings.Host = "localhost"
	settings.Port = port
	client, err := iec61850.NewClient(settings)
	if err != nil {
		t.Fatalf("client connect: %v", err)
	}
	defer client.Close()

	opErr := client.ControlByControlModel(pcControlRef, iec61850.CONTROL_MODEL_DIRECT_NORMAL, iec61850.NewControlObjectParam(true))
	stVal, _ := client.ReadBool(pcStValRef, iec61850.ST)
	return opErr, stVal
}

func TestPerformCheckHandler_DenyProducesOperateError(t *testing.T) {
	deny := func(_ *iec61850.ModelNode, action *iec61850.ControlAction, _ *iec61850.MmsValue, _ bool, _ bool) iec61850.CheckHandlerResult {
		action.SetAddCause(iec61850.ADD_CAUSE_TIME_LIMIT_OVER)
		return iec61850.CONTROL_OBJECT_ACCESS_DENIED
	}

	opErr, stVal := operateWithCheck(t, 10298, deny)
	t.Logf("deny: operate err=%v, stVal=%v", opErr, stVal)

	if opErr == nil {
		t.Fatal("expected the client Operate to return an error when the perform-check handler denies, got nil")
	}
	if stVal {
		t.Errorf("a denied Operate must not change the value, but stVal=true")
	}
}

func TestPerformCheckHandler_AcceptSucceeds(t *testing.T) {
	accept := func(_ *iec61850.ModelNode, _ *iec61850.ControlAction, _ *iec61850.MmsValue, _ bool, _ bool) iec61850.CheckHandlerResult {
		return iec61850.CONTROL_ACCEPTED
	}

	opErr, stVal := operateWithCheck(t, 10299, accept)
	t.Logf("accept: operate err=%v, stVal=%v", opErr, stVal)

	if opErr != nil {
		t.Fatalf("expected an accepted Operate to succeed, got err=%v", opErr)
	}
}
