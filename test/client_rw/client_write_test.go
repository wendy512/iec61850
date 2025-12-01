package client_rw

import (
	"github.com/boeboe/iec61850"
	"github.com/boeboe/iec61850/test"
	"testing"
)

const (
	OutVarObjectRef = "ied1Inverter/ZINV1.OutVarSet.setMag.f"
)

func TestWrite(t *testing.T) {
	client := test.CreateClient(t)
	defer test.CloseClient(client)

	if err := client.Write(OutVarObjectRef, iec61850.SP, 100); err != nil {
		t.Fatalf("write %s error %v\n", OutVarObjectRef, err)
	}
}
