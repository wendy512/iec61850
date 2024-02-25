package test

import (
	"github.com/wendy512/iec61850/scl_xml"
	"testing"
)

func TestLoadIcdXml(t *testing.T) {
	scl, err := scl_xml.GetSCL("test.icd")
	if err != nil {
		t.Error(err)
	}
	scl.Print()
}
