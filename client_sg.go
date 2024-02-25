package iec61850

// #include <iec61850_client.h>
import "C"
import (
	"fmt"
	"github.com/spf13/cast"
	"unsafe"
)

const (
	ActDA  = "%s/%s.SGCB.ActSG"
	EditDA = "%s/%s.SGCB.EditSG"
	CnfDA  = "%s/%s.SGCB.CnfEdit"
)

type SettingGroup struct {
	NumOfSG int
	ActSG   int
	EditSG  int
	CnfEdit bool
}

// WriteSG 写入SettingGroup
func (c *Client) WriteSG(ld, ln, objectRef string, fc FC, actSG int, value interface{}) error {
	// Set active setting group
	if err := c.Write(fmt.Sprintf(ActDA, ld, ln), SP, actSG); err != nil {
		return err
	}

	// Set edit setting group
	if err := c.Write(fmt.Sprintf(EditDA, ld, ln), SP, actSG); err != nil {
		return err
	}

	// Change a setting group value
	if err := c.Write(objectRef, fc, value); err != nil {
		return err
	}

	// Confirm new setting group values
	if err := c.Write(fmt.Sprintf(CnfDA, ld, ln), SP, true); err != nil {
		return err
	}
	return nil
}

// GetSG 获取SettingGroup
func (c *Client) GetSG(objectRef string) (*SettingGroup, error) {
	var clientError C.IedClientError
	cObjectRef := C.CString(objectRef)
	defer C.free(unsafe.Pointer(cObjectRef))

	// 获取类型
	sgcbVarSpec := C.IedConnection_getVariableSpecification(c.conn, &clientError, cObjectRef, C.FunctionalConstraint(SP))
	if err := GetIedClientError(clientError); err != nil {
		return nil, err
	}
	defer C.MmsVariableSpecification_destroy(sgcbVarSpec)

	// Read SGCB
	sgcbVal := C.IedConnection_readObject(c.conn, &clientError, cObjectRef, C.FunctionalConstraint(SP))
	if err := GetIedClientError(clientError); err != nil {
		return nil, err
	}
	defer C.MmsValue_delete(sgcbVal)

	numOfSGValue := c.getSubElementValue(sgcbVal, sgcbVarSpec, "NumOfSG")
	actSGValue := c.getSubElementValue(sgcbVal, sgcbVarSpec, "ActSG")
	editSGValue := c.getSubElementValue(sgcbVal, sgcbVarSpec, "EditSG")
	cnfEditValue := c.getSubElementValue(sgcbVal, sgcbVarSpec, "CnfEdit")

	sg := &SettingGroup{
		NumOfSG: cast.ToInt(numOfSGValue),
		ActSG:   cast.ToInt(actSGValue),
		EditSG:  cast.ToInt(editSGValue),
		CnfEdit: cast.ToBool(cnfEditValue),
	}
	return sg, nil
}
