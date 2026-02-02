package goose_publisher

import (
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/wendy512/iec61850"
)

func TestGoosePublisher(t *testing.T) {
	publisher, err := iec61850.NewGoosePublisher(iec61850.GoosePublisherConf{
		InterfaceID: "eth0",
		AppID:       1000,
		DstAddr: [6]uint8{
			0x01,
			0x0c,
			0xcd,
			0x01,
			0x00,
			0x01,
		},
		VlanID:       0,
		VlanPriority: 4,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer publisher.Close()

	publisher.SetGoCbRef("simpleIOGenericIO/LLN0$GO$gcbAnalogValues")
	publisher.SetDataSetRef("simpleIOGenericIO/LLN0$AnalogValues")
	publisher.SetConfRev(1)
	publisher.SetTimeAllowedToLive(500)

	lVal := iec61850.NewLinkedListValue()
	defer lVal.Destroy()
	if err := lVal.Add(&iec61850.MmsValue{
		Type:  iec61850.Int64,
		Value: time.Now().UnixMilli(),
	}); err != nil {
		t.Fatal(err)
	}
	if err := lVal.Add(&iec61850.MmsValue{
		Type:  iec61850.Float,
		Value: 233.3,
	}); err != nil {
		t.Fatal(err)
	}
	t.Logf("linked list size is %d\n", lVal.Size())
	ticker := time.NewTicker(time.Second * 2)
	defer ticker.Stop()
	go func() {
		var i uint64
		for range ticker.C {
			if err := publisher.Publish(lVal); err != nil {
				t.Error(err)
				return
			} else {
				i++
				t.Logf("send [%d] goose message successfully", i)
			}
		}
	}()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	t.Logf("Goose  Publisher Close\n")
}
