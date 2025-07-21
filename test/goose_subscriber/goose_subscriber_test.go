package goose_subscriber

import (
	"encoding/xml"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/wendy512/iec61850"
)

func printReportValue(report *iec61850.GooseReport) {
	fmt.Printf("[appID: %d, goID: %s]\n", report.GetAppID(), report.GetGoID())
	fmt.Printf("[vlanID: %d, vlanPriority: %d]\n", report.GetVlanID(), report.GetVlanPriority())
	fmt.Printf("[dstMac: %#v, srcMac: %#v]\n", report.GetDstMac(), report.GetSrcMac())
	fmt.Printf("[goCbRef: %s, dataSetName: %s]\n", report.GetGoCbRef(), report.GetDataSetName())
	fmt.Printf("[timestamp: %d, time allowed to live: %d]\n", report.GetTimestamp(), report.GetTimeAllowedToLive())
	fmt.Printf("[stNum: %d, sqNum: %d]\n", report.GetStNum(), report.GetSqNum())

	if v, err := report.GetDataSetValues(); err != nil {
		fmt.Println("to GO value error: ", err)
	} else {
		str, _ := xml.MarshalIndent(v, "", "  ")
		fmt.Println(string(str))
	}
}

func TestGooseSubscriber(t *testing.T) {
	gooseReceiver := iec61850.NewGooseReceiver()
	defer gooseReceiver.Stop().Destroy()

	subscriber1 := iec61850.NewGooseSubscriber(iec61850.SubscriberConf{
		InterfaceID:   "eth0",
		AppID:         1000,
		DstMacAddr:    [6]uint8{0x01, 0x0c, 0xcd, 0x01, 0x00, 0x01},
		Subscriber:    "simpleIOGenericIO/LLN0$GO$gcbAnalogValues",
		ReportHandler: printReportValue,
	})

	if !gooseReceiver.
		AddSubscriber(subscriber1).
		Start().
		IsRunning() {
		t.Fatal("can't start goose receiver")
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	t.Logf("Goose Subscriber Close\n")
}
