package test

import (
	"github.com/wendy512/iec61850"
	"testing"
)

func TestControl(t *testing.T) {
	client := createClient(t)
	objectRef := "simpleIOGenericIO/GGIO1.SPCSO1"
	value := true

	if err := client.ControlForDirectWithNormalSecurity(objectRef, value); err != nil {
		t.Errorf("[direct-with-normal-security] %s object error %v\n", objectRef, err)
		t.FailNow()
	}
	doRead(t, client, objectRef+".stVal", iec61850.ST)

	objectRef = "simpleIOGenericIO/GGIO1.SPCSO2"
	if err := client.ControlForSboWithNormalSecurity(objectRef, value); err != nil {
		t.Errorf("[sbo-with-normal-security] %s object error %v\n", objectRef, err)
		t.FailNow()
	}
	doRead(t, client, objectRef+".stVal", iec61850.ST)

	objectRef = "simpleIOGenericIO/GGIO1.SPCSO3"
	if err := client.ControlForDirectWithEnhancedSecurity(objectRef, value); err != nil {
		t.Errorf("[direct-with-enhanced-security] %s object error %v\n", objectRef, err)
		t.FailNow()
	}
	doRead(t, client, objectRef+".stVal", iec61850.ST)

	objectRef = "simpleIOGenericIO/GGIO1.SPCSO4"
	if err := client.ControlForSboWithEnhancedSecurity(objectRef, value); err != nil {
		t.Errorf("[sbo-with-enhanced-security] %s object error %v\n", objectRef, err)
		t.FailNow()
	}
	doRead(t, client, objectRef+".stVal", iec61850.ST)
}
