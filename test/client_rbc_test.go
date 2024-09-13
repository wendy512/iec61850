package test

import (
	"fmt"
	"testing"

	"github.com/morris-kelly/iec61850"
)

func TestRBC(t *testing.T) {
	client, err := iec61850.NewClient(iec61850.NewSettings())
	if err != nil {
		t.Fatal(err)
	}
	read, err := client.Read("CL10002ALD0/DevAlmGGIO1.Beh.stVal", iec61850.ST)
	if err != nil {
		return
	}
	fmt.Println(read)
	rbcRef := "CL10002ALD0/LLN0.BR.brcbRelayEna02"
	values, err := client.ReadRbcValues(rbcRef)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(values)
	err = client.SetRbcValues(rbcRef, iec61850.ClientReportControlBlock{
		Ena:    true,
		IntgPd: 500,
		OptFlds: iec61850.OptFlds{
			SequenceNumber:     true,
			TimeOfEntry:        true,
			ReasonForInclusion: true,
			DataSetName:        true,
			DataReference:      true,
			BufferOverflow:     true,
			EntryID:            true,
			ConfigRevision:     true,
		},
		TrgOps: iec61850.TrgOps{
			DataChange:            true,
			QualityChange:         true,
			DataUpdate:            true,
			TriggeredPeriodically: true,
			Gi:                    true,
			Transient:             false,
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
	values, err = client.ReadRbcValues(rbcRef)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(values)
}
