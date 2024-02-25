package iec61850

// #include <iec61850_client.h>
import "C"
import "unsafe"

// ControlForDirectWithNormalSecurity 控制模式 1[direct-with-normal-security]
func (c *Client) ControlForDirectWithNormalSecurity(objectRef string, value bool) error {
	return c.control(objectRef, value, false, true, false)
}

// ControlForSboWithNormalSecurity 控制模式 2[sbo-with-normal-security]
func (c *Client) ControlForSboWithNormalSecurity(objectRef string, value bool) error {
	return c.control(objectRef, value, true, false, false)
}

// ControlForDirectWithEnhancedSecurity 控制模式 3[direct-with-enhanced-security]
func (c *Client) ControlForDirectWithEnhancedSecurity(objectRef string, value bool) error {
	return c.control(objectRef, value, false, false, true)
}

// ControlForSboWithEnhancedSecurity 控制模式 4[sbo-with-enhanced-security]
func (c *Client) ControlForSboWithEnhancedSecurity(objectRef string, value bool) error {
	return c.control(objectRef, value, true, false, true)
}

func (c *Client) control(objectRef string, value, _select, direct, enhanced bool) error {
	cObjectRef := C.CString(objectRef)
	defer C.free(unsafe.Pointer(cObjectRef))

	control := C.ControlObjectClient_create(cObjectRef, c.conn)
	if control == nil {
		return CreateControlObjectClientFail
	}

	// Select before operate
	defer C.ControlObjectClient_destroy(control)
	if _select && !enhanced && !bool(C.ControlObjectClient_select(control)) {
		return ControlSelectFail
	}

	// Direct control with enhanced security
	if enhanced {
		C.ControlObjectClient_setCommandTerminationHandler(control, nil, nil)
	}
	ctlVal := C.MmsValue_newBoolean(C.bool(value))
	defer C.MmsValue_delete(ctlVal)

	// Select before operate with enhanced security
	if _select && enhanced && !bool(C.ControlObjectClient_selectWithValue(control, ctlVal)) {
		return ControlSelectFail
	}

	// Direct control
	if direct {
		C.ControlObjectClient_setOrigin(control, nil, 3)
	}

	C.ControlObjectClient_operate(control, ctlVal, 0)
	return nil
}
