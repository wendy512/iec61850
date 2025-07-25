package sv_publisher

import (
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/wendy512/iec61850"
)

func TestSvPublisher(t *testing.T) {
	publisher1 := iec61850.NewSVPublisher(iec61850.SvPublisherConf{
		EtherName:    "eth0",
		AppID:        12401,
		VlanID:       1,
		VlanPriority: 4,
	})
	defer publisher1.Destroy()

	publisher2 := iec61850.NewSVPublisher(iec61850.SvPublisherConf{
		EtherName:    "eth0",
		AppID:        12403,
		VlanID:       2,
		VlanPriority: 5,
	})
	defer publisher2.Destroy()

	ticker1 := time.NewTicker(time.Second * 2)
	ticker2 := time.NewTicker(time.Second * 3)
	defer ticker1.Stop()
	defer ticker2.Stop()

	go func() {
		fVal1 := 0.1
		fVal2 := 0.2
		timestamp3 := iec61850.NewTimestamp()

		asdu1 := publisher1.AddSvASDU("svpub1", "svpub1", 1)
		asdu2 := publisher1.AddSvASDU("svpub2", "svpub2", 1)

		floatIndex1 := asdu1.AddFloat()
		floatIndex2 := asdu1.AddFloat()
		timeIndex3 := asdu1.AddTimestamp()

		floatIndex4 := asdu2.AddFloat()
		floatIndex5 := asdu2.AddFloat()
		timeIndex6 := asdu2.AddTimestamp()
		publisher1.SetupComplete()

		for range ticker1.C {
			asdu1.SetFloat(floatIndex1, float32(fVal1))
			asdu1.SetFloat(floatIndex2, float32(fVal2))
			asdu1.SetTimestamp(timeIndex3, timestamp3)
			asdu1.IncreaseSmpCnt()

			asdu2.SetFloat(floatIndex4, float32(fVal1))
			asdu2.SetFloat(floatIndex5, float32(fVal2))
			asdu2.SetTimestamp(timeIndex6, timestamp3)
			asdu2.IncreaseSmpCnt()

			publisher1.Publish()
			fVal1 += 0.1
			fVal2 += 0.3
			timestamp3.SetTime(time.Now())
		}
	}()

	go func() {
		fVal1 := 0.3
		fVal2 := 0.9
		timestamp3 := iec61850.NewTimestamp()
		asdu1 := publisher2.AddSvASDU("svpub3", "svpub3", 1)
		asdu2 := publisher2.AddSvASDU("svpub4", "svpub4", 1)

		floatIndex1 := asdu1.AddFloat()
		floatIndex2 := asdu1.AddFloat()
		timeIndex3 := asdu1.AddTimestamp()

		floatIndex4 := asdu2.AddFloat()
		floatIndex5 := asdu2.AddFloat()
		timeIndex6 := asdu2.AddTimestamp()
		publisher2.SetupComplete()

		for range ticker1.C {
			asdu1.SetFloat(floatIndex1, float32(fVal1))
			asdu1.SetFloat(floatIndex2, float32(fVal2))
			asdu1.SetTimestamp(timeIndex3, timestamp3)
			asdu1.IncreaseSmpCnt()

			asdu2.SetFloat(floatIndex4, float32(fVal1))
			asdu2.SetFloat(floatIndex5, float32(fVal2))
			asdu2.SetTimestamp(timeIndex6, timestamp3)
			asdu2.IncreaseSmpCnt()

			publisher2.Publish()
			fVal1 += 0.1
			fVal2 += 0.3
			timestamp3.SetTime(time.Now())
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	t.Logf("All Sv Publisher Close\n")
}
