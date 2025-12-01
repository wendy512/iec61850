package client_rw

import (
	"encoding/json"
	"fmt"
	"github.com/boeboe/iec61850"
	"github.com/boeboe/iec61850/test"
	"testing"
)

const (
	AnIn1ObjectRef = "simpleIOGenericIO/GGIO1.AnIn1.mag.f"
	Ind1ObjectRef  = "simpleIOGenericIO/GGIO1.Ind1.stVal"
)

func TestRead(t *testing.T) {
	client := test.CreateClient(t)
	defer test.CloseClient(client)

	test.DoRead(t, client, AnIn1ObjectRef, iec61850.MX)
	test.DoRead(t, client, Ind1ObjectRef, iec61850.ST)
}

func TestReadFloat(t *testing.T) {
	client := test.CreateClient(t)
	defer test.CloseClient(client)
	objectRef := AnIn1ObjectRef

	value, err := client.ReadFloat(objectRef, iec61850.MX)
	if err != nil {
		t.Fatalf("read %s object error %v\n", objectRef, err)
	}
	t.Logf("read %s float value -> %v\n", objectRef, value)
}

func TestReadBool(t *testing.T) {
	client := test.CreateClient(t)
	defer test.CloseClient(client)

	objectRef := Ind1ObjectRef
	value, err := client.ReadBool(objectRef, iec61850.ST)
	if err != nil {
		t.Fatalf("read %s object error %v\n", objectRef, err)
	}
	t.Logf("read %s bool value -> %v\n", objectRef, value)
}

func TestGetLogicalDeviceList(t *testing.T) {
	client := test.CreateClient(t)
	defer test.CloseClient(client)

	deviceList := client.GetLogicalDeviceList()
	marshal, err := json.Marshal(deviceList)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(marshal))
}
