package client_control

import (
	"github.com/wendy512/iec61850"
	"github.com/wendy512/iec61850/test"
	"testing"
)

const DefValue = false

func TestControl(t *testing.T) {
	client := test.CreateClient(t)
	defer test.CloseClient(client)

	objectRef := "simpleIOGenericIO/GGIO1.SPCSO1"

	param := iec61850.NewControlObjectParam(DefValue)
	param.OrIdent = "test"
	if err := client.ControlByControlModel(objectRef, iec61850.CONTROL_MODEL_DIRECT_NORMAL, param); err != nil {
		t.Errorf("[direct-with-normal-security] %s object error %v\n", objectRef, err)
		t.FailNow()
	}
	test.DoRead(t, client, objectRef+".stVal", iec61850.ST)

	objectRef = "simpleIOGenericIO/GGIO1.SPCSO2"
	if err := client.ControlForSboWithNormalSecurity(objectRef, DefValue); err != nil {
		t.Errorf("[sbo-with-normal-security] %s object error %v\n", objectRef, err)
		t.FailNow()
	}
	test.DoRead(t, client, objectRef+".stVal", iec61850.ST)

	objectRef = "simpleIOGenericIO/GGIO1.SPCSO3"
	if err := client.ControlForDirectWithEnhancedSecurity(objectRef, DefValue); err != nil {
		t.Errorf("[direct-with-enhanced-security] %s object error %v\n", objectRef, err)
		t.FailNow()
	}
	test.DoRead(t, client, objectRef+".stVal", iec61850.ST)

	objectRef = "simpleIOGenericIO/GGIO1.SPCSO4"
	if err := client.ControlForSboWithEnhancedSecurity(objectRef, DefValue); err != nil {
		t.Errorf("[sbo-with-enhanced-security] %s object error %v\n", objectRef, err)
		t.FailNow()
	}
	test.DoRead(t, client, objectRef+".stVal", iec61850.ST)
}
