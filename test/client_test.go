package test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/morris-kelly/iec61850"
)

const (
	AnIn1ObjectRef = "simpleIOGenericIO/GGIO1.AnIn1.mag.f"
	Ind1ObjectRef  = "simpleIOGenericIO/GGIO1.Ind1.stVal"
)

func TestRead(t *testing.T) {
	client := createClient(t)
	doRead(t, client, AnIn1ObjectRef, iec61850.MX)
	doRead(t, client, Ind1ObjectRef, iec61850.ST)
}

func TestReadFloat(t *testing.T) {
	client := createClient(t)
	objectRef := AnIn1ObjectRef

	value, err := client.ReadFloat(objectRef, iec61850.MX)
	if err != nil {
		t.Errorf("read %s object error %v\n", objectRef, err)
		t.FailNow()
	}
	t.Logf("read %s float value -> %v", objectRef, value)
}

func TestReadBool(t *testing.T) {
	client := createClient(t)
	objectRef := Ind1ObjectRef

	value, err := client.ReadBool(objectRef, iec61850.ST)
	if err != nil {
		t.Errorf("read %s object error %v\n", objectRef, err)
		t.FailNow()
	}
	t.Logf("read %s bool value -> %v", objectRef, value)
}

func createClient(t *testing.T) *iec61850.Client {
	client, err := iec61850.NewClient(iec61850.NewSettings())
	if err != nil {
		t.Errorf("create client error %v\n", err)
		t.FailNow()
	}
	return client
}

func doRead(t *testing.T, client *iec61850.Client, objectRef string, fc iec61850.FC) {
	value, err := client.Read(objectRef, fc)
	if err != nil {
		t.Errorf("read %s object error %v\n", objectRef, err)
		t.FailNow()
	}
	t.Logf("read %s value -> %v", objectRef, value)
}

func TestGetLogicalDeviceList(t *testing.T) {
	client, err := iec61850.NewClient(&iec61850.Settings{
		Host:           "127.0.0.1",
		Port:           10086,
		ConnectTimeout: 10000,
		RequestTimeout: 10000,
	})
	if err != nil {
		panic(err)
	}
	deviceList := client.GetLogicalDeviceList()
	marshal, err := json.Marshal(deviceList)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(marshal))
}
