package scl

import (
	"testing"

	"github.com/marrasen/iec61850/scl_xml"
)

func TestLoadIcdXml(t *testing.T) {
	scl, err := scl_xml.GetSCL("test.icd")
	if err != nil {
		t.Error(err)
	}
	scl.Print()
}
