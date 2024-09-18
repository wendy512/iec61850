package test

import (
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/wendy512/iec61850"
)

func TestCreateServerFromConfigFile(t *testing.T) {
	model, err := iec61850.CreateModelFromConfigFileEx("model.cfg")
	if err != nil {
		t.Errorf("create model error %v\n", err)
		t.FailNow()
	}

	server := iec61850.NewServer(model)
	server.Start(102)
	defer server.Destroy()
	t.Logf("Server start up")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	<-sig
	t.Logf("Server stop")
	server.Stop()
}
