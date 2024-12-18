package tls_client

import (
	"github.com/wendy512/iec61850"
	"github.com/wendy512/iec61850/test"
	"testing"
)

const (
	AnIn1ObjectRef = "simpleIOGenericIO/GGIO1.AnIn1.mag.f"
	Ind1ObjectRef  = "simpleIOGenericIO/GGIO1.Ind1.stVal"
)

func TestReadForTls(t *testing.T) {
	client := test.CreateTlsClient(t)
	defer test.CloseClient(client)

	test.DoRead(t, client, AnIn1ObjectRef, iec61850.MX)
	test.DoRead(t, client, Ind1ObjectRef, iec61850.ST)

	objectRef := "simpleIOGenericIO/GGIO1.SPCSO1"
	param := iec61850.NewControlObjectParam(true)
	param.OrIdent = "test"
	if err := client.ControlByControlModel(objectRef, iec61850.CONTROL_MODEL_DIRECT_NORMAL, param); err != nil {
		t.Errorf("[direct-with-normal-security] %s object error %v\n", objectRef, err)
		t.FailNow()
	}
	test.DoRead(t, client, objectRef+".stVal", iec61850.ST)
}
