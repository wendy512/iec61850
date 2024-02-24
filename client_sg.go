package iec61850

// #include <iec61850_client.h>
import "C"
import (
	"github.com/spf13/cast"
	"unsafe"
)

const (
	ActDA  = "/LLN0.SGCB.ActSG"
	EditDA = "/LLN0.SGCB.EditSG"
	CnfDA  = "/LLN0.SGCB.CnfEdit"
)

type SettingGroup struct {
	NumOfSG int
	ActSG   int
	EditSG  int
	CnfEdit bool
}

func (c *Client) WriteSG(ld, objectRef string, fc FC, actSG int, value interface{}) error {
	// Set active setting group
	if err := c.Write(ld+ActDA, SP, actSG); err != nil {
		return err
	}

	// Set edit setting group
	if err := c.Write(ld+EditDA, SP, actSG); err != nil {
		return err
	}

	// Change a setting group value
	if err := c.Write(objectRef, fc, value); err != nil {
		return err
	}

	// Confirm new setting group values
	if err := c.Write(ld+CnfDA, SP, true); err != nil {
		return err
	}
	return nil
}

func (c *Client) GetSG(ld string) (*SettingGroup, error) {
	var clientError C.IedClientError
	cObjectRef := C.CString(ld + "/LLN0.SGCB")
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

	numOfSGValue := c.getSGSubElementValue(sgcbVal, sgcbVarSpec, "NumOfSG")
	actSGValue := c.getSGSubElementValue(sgcbVal, sgcbVarSpec, "ActSG")
	editSGValue := c.getSGSubElementValue(sgcbVal, sgcbVarSpec, "EditSG")
	cnfEditValue := c.getSGSubElementValue(sgcbVal, sgcbVarSpec, "CnfEdit")

	sg := &SettingGroup{
		NumOfSG: cast.ToInt(numOfSGValue),
		ActSG:   cast.ToInt(actSGValue),
		EditSG:  cast.ToInt(editSGValue),
		CnfEdit: cast.ToBool(cnfEditValue),
	}
	return sg, nil
}

func (c *Client) getSGSubElementValue(sgcbVal *C.MmsValue, sgcbVarSpec *C.MmsVariableSpecification, name string) interface{} {
	mmsPath := C.CString(name)
	defer C.free(unsafe.Pointer(mmsPath))
	mmsValue := C.MmsValue_getSubElement(sgcbVal, sgcbVarSpec, mmsPath)
	defer C.MmsValue_delete(mmsValue)
	return c.toGoValue(mmsValue, MmsType(C.MmsValue_getType(mmsValue)))
}
