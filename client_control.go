package iec61850

// #include <iec61850_client.h>
import "C"
import "unsafe"

type commandTerminationHandler func(parameter unsafe.Pointer, control C.ControlObjectClient)

// DirectControl 直接控制
func (c *Client) DirectControl(objectRef string, value bool) error {
	return c.control(objectRef, value, false, true, false, nil)
}

// DirectControlSelect 直接控制之前select
func (c *Client) DirectControlSelect(objectRef string, value bool) error {
	return c.control(objectRef, value, true, false, false, nil)
}

// DirectControlWithEnhancedSecurity 增强安全直接控制
func (c *Client) DirectControlWithEnhancedSecurity(objectRef string, value bool) error {
	return c.control(objectRef, value, false, false, true, func(void unsafe.Pointer, control C.ControlObjectClient) {})
}

// DirectControlSelectWithEnhancedSecurity 增强安全直接控制之前select
func (c *Client) DirectControlSelectWithEnhancedSecurity(objectRef string, value bool) error {
	return c.control(objectRef, value, true, false, true, nil)
}

// control for FC=CO
func (c *Client) control(objectRef string, value, _select, direct, enhanced bool, handler commandTerminationHandler) error {
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
		C.ControlObjectClient_setCommandTerminationHandler(control, C.CommandTerminationHandler(handler), nil)
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
	if bool(C.ControlObjectClient_operate(control, ctlVal, 0)) {
		return nil
	} else {
		return ControlObjectFail
	}
}
