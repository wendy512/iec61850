//go:build linux && amd64

package sv_subscriber

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"testing"

	"github.com/boeboe/iec61850"
)

func createPrintReporter(sep string) iec61850.SvReportHandler {
	return func(report *iec61850.SvReport) {
		asdu := report.ReceiverASDU
		fmt.Println(strings.Repeat(sep, 20))
		fmt.Printf("confRev: %d; dataSet: %s; dataSize: %d\n", asdu.GetConfRev(), asdu.GetDatSet(), asdu.GetDataSize())
		fmt.Printf("timeAsMs: %d; timeAsNs: %d\n", asdu.GetRefrTmAsMs(), asdu.GetRefrTmAsNs())
		fmt.Printf("DATA[0]: %f; DATA[4]: %f\n", asdu.GetFloat32(0), asdu.GetFloat32(4))
		fmt.Printf("Time %v; ClockFailure: %v; ClockNotSync: %v; IsLeapSecondKnown: %v\n",
			asdu.GetTimestamp(8).GetTime(),
			asdu.GetTimestamp(8).HasClockFailure(),
			asdu.GetTimestamp(8).IsClockNotSynchronized(),
			asdu.GetTimestamp(8).IsLeapSecondKnown(),
		)
		fmt.Println(strings.Repeat(sep, 20))
		fmt.Println()
	}
}

func TestSvSubscriber(t *testing.T) {
	receiver := iec61850.NewSvReceiver(iec61850.SvReceiverConf{
		InterfaceID: "eth0",
	})
	defer receiver.Stop().Destroy()
	var sub1 = iec61850.NewSvSubscriber(iec61850.SvSubscriberConf{
		AppID:   12401,
		Handler: createPrintReporter("--"),
	})
	var sub2 = iec61850.NewSvSubscriber(iec61850.SvSubscriberConf{
		AppID:   12403,
		Handler: createPrintReporter("**"),
	})

	if !receiver.
		AddSubscriber(sub1).
		AddSubscriber(sub2).
		Start().
		IsRunning() {
		t.Fatal("can not startup sv receiver")
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	t.Logf("All SvSubscriber Close\n")
}
