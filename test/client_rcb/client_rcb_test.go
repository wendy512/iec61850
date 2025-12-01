package client_rcb

import (
	"github.com/boeboe/iec61850"
	"github.com/boeboe/iec61850/test"
	"log"
	"testing"
)

func TestRBC(t *testing.T) {
	client := test.CreateClient(t)
	defer test.CloseClient(client)

	objectRef := "simpleIOGenericIO/GGIO1.SPCSO1.stVal"
	value, err := client.Read(objectRef, iec61850.ST)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("%s -> %v\n", objectRef, value)

	rbcRef := "simpleIOGenericIO/LLN0.RP.EventsRCB01"
	rcbValue, err := client.GetRCBValues(rbcRef)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("write before %s -> %#v\n", rbcRef, rcbValue)

	err = client.SetRCBValues(rbcRef, iec61850.ClientReportControlBlock{
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
		t.Fatal(err)
	}

	rcbValue, err = client.GetRCBValues(rbcRef)
	if err != nil {
		t.Error(err)
		return
	}
	log.Printf("write after %s -> %#v\n", rbcRef, rcbValue)
}
