package iec61850

// #include <iec61850_client.h>
import "C"
import "unsafe"

type ControlObjectParam struct {
	CtlVal      bool
	OrIdent     string
	OrCat       int
	Test        bool
	Check       bool
	OperateTime uint64
}

type ControlObjectParamAPC struct {
	CtlVal      float32
	OrIdent     string
	OrCat       int
	Test        bool
	Check       bool
	OperateTime uint64
}

type ControlObjectParamINC struct {
	CtlVal      int
	OrIdent     string
	OrCat       int
	Test        bool
	Check       bool
	OperateTime uint64
}

func NewControlObjectParam(ctlVal bool) *ControlObjectParam {
	return &ControlObjectParam{
		CtlVal:      ctlVal,
		OrIdent:     "",
		OrCat:       0,
		Test:        false,
		Check:       false,
		OperateTime: 0,
	}
}

func NewControlObjectParamAPC(ctlVal float32) *ControlObjectParamAPC {
	return &ControlObjectParamAPC{
		CtlVal:      ctlVal,
		OrIdent:     "",
		OrCat:       0,
		Test:        false,
		Check:       false,
		OperateTime: 0,
	}
}

func NewControlObjectParamINC(ctlVal int) *ControlObjectParamINC {
	return &ControlObjectParamINC{
		CtlVal:      ctlVal,
		OrIdent:     "",
		OrCat:       0,
		Test:        false,
		Check:       false,
		OperateTime: 0,
	}
}

// ControlForDirectWithNormalSecurity 控制模式 1[direct-with-normal-security]
func (c *Client) ControlForDirectWithNormalSecurity(objectRef string, ctlVal bool) error {
	return c.ControlByControlModel(objectRef, CONTROL_MODEL_DIRECT_NORMAL, NewControlObjectParam(ctlVal))
}

func (c *Client) ControlByControlModelINC(objectRef string, controlModel ControlModel, param *ControlObjectParamINC) error {
	cObjectRef := C.CString(objectRef)
	defer C.free(unsafe.Pointer(cObjectRef))

	control := C.ControlObjectClient_create(cObjectRef, c.conn)
	if control == nil {
		return CreateControlObjectClientFail
	}

	ctlVal := C.MmsValue_newIntegerFromInt32(C.int(param.CtlVal))
	defer C.MmsValue_delete(ctlVal)

	switch controlModel {
	case CONTROL_MODEL_SBO_NORMAL:
		if !bool(C.ControlObjectClient_select(control)) {
			return ControlSelectFail
		}
	case CONTROL_MODEL_DIRECT_ENHANCED:
		C.ControlObjectClient_setCommandTerminationHandler(control, nil, nil)
	case CONTROL_MODEL_SBO_ENHANCED:
		C.ControlObjectClient_setCommandTerminationHandler(control, nil, nil)
		if !bool(C.ControlObjectClient_selectWithValue(control, ctlVal)) {
			return ControlSelectFail
		}
	}

	var cOrIdent *C.char
	if param.OrIdent != "" {
		cOrIdent = C.CString(param.OrIdent)
		defer C.free(unsafe.Pointer(cOrIdent))
	}

	C.ControlObjectClient_setControlModel(control, C.ControlModel(controlModel))
	C.ControlObjectClient_setOrigin(control, cOrIdent, C.int(param.OrCat))
	C.ControlObjectClient_setInterlockCheck(control, C.bool(param.Check))
	C.ControlObjectClient_setSynchroCheck(control, C.bool(param.Check))
	C.ControlObjectClient_setTestMode(control, C.bool(param.Test))

	if !bool(C.ControlObjectClient_operate(control, ctlVal, 0)) {
		return ControlObjectFail
	}
	return nil
}

func (c *Client) ControlByControlModelAPC(objectRef string, controlModel ControlModel, param *ControlObjectParamAPC) error {
	cObjectRef := C.CString(objectRef)
	defer C.free(unsafe.Pointer(cObjectRef))

	control := C.ControlObjectClient_create(cObjectRef, c.conn)
	if control == nil {
		return CreateControlObjectClientFail
	}

	ctlVal := C.MmsValue_newFloat(C.float(param.CtlVal))
	defer C.MmsValue_delete(ctlVal)

	switch controlModel {
	case CONTROL_MODEL_SBO_NORMAL:
		if !bool(C.ControlObjectClient_select(control)) {
			return ControlSelectFail
		}
	case CONTROL_MODEL_DIRECT_ENHANCED:
		C.ControlObjectClient_setCommandTerminationHandler(control, nil, nil)
	case CONTROL_MODEL_SBO_ENHANCED:
		C.ControlObjectClient_setCommandTerminationHandler(control, nil, nil)
		if !bool(C.ControlObjectClient_selectWithValue(control, ctlVal)) {
			return ControlSelectFail
		}
	}

	var cOrIdent *C.char
	if param.OrIdent != "" {
		cOrIdent = C.CString(param.OrIdent)
		defer C.free(unsafe.Pointer(cOrIdent))
	}

	C.ControlObjectClient_setControlModel(control, C.ControlModel(controlModel))
	C.ControlObjectClient_setOrigin(control, cOrIdent, C.int(param.OrCat))
	C.ControlObjectClient_setInterlockCheck(control, C.bool(param.Check))
	C.ControlObjectClient_setSynchroCheck(control, C.bool(param.Check))
	C.ControlObjectClient_setTestMode(control, C.bool(param.Test))

	if !bool(C.ControlObjectClient_operate(control, ctlVal, 0)) {
		return ControlObjectFail
	}
	return nil
}

func (c *Client) ControlByControlModel(objectRef string, controlModel ControlModel, param *ControlObjectParam) error {
	cObjectRef := C.CString(objectRef)
	defer C.free(unsafe.Pointer(cObjectRef))

	control := C.ControlObjectClient_create(cObjectRef, c.conn)
	if control == nil {
		return CreateControlObjectClientFail
	}

	defer C.ControlObjectClient_destroy(control)
	ctlVal := C.MmsValue_newBoolean(C.bool(param.CtlVal))
	defer C.MmsValue_delete(ctlVal)

	// Select before operate
	switch controlModel {
	case CONTROL_MODEL_SBO_NORMAL:
		if !bool(C.ControlObjectClient_select(control)) {
			return ControlSelectFail
		}
	case CONTROL_MODEL_DIRECT_ENHANCED:
		C.ControlObjectClient_setCommandTerminationHandler(control, nil, nil)
	case CONTROL_MODEL_SBO_ENHANCED:
		C.ControlObjectClient_setCommandTerminationHandler(control, nil, nil)
		if !bool(C.ControlObjectClient_selectWithValue(control, ctlVal)) {
			return ControlSelectFail
		}
	}

	var cOrIdent *C.char
	if param.OrIdent != "" {
		cOrIdent = C.CString(param.OrIdent)
		defer C.free(unsafe.Pointer(cOrIdent))
	}

	C.ControlObjectClient_setControlModel(control, C.ControlModel(controlModel))
	C.ControlObjectClient_setOrigin(control, cOrIdent, C.int(param.OrCat))
	C.ControlObjectClient_setInterlockCheck(control, C.bool(param.Check))
	C.ControlObjectClient_setSynchroCheck(control, C.bool(param.Check))
	C.ControlObjectClient_setTestMode(control, C.bool(param.Test))

	if !bool(C.ControlObjectClient_operate(control, ctlVal, C.uint64_t(param.OperateTime))) {
		return ControlObjectFail
	}
	return nil
}

// ControlForSboWithNormalSecurity 控制模式 2[sbo-with-normal-security]
func (c *Client) ControlForSboWithNormalSecurity(objectRef string, value bool) error {
	return c.ControlByControlModel(objectRef, CONTROL_MODEL_SBO_NORMAL, NewControlObjectParam(value))
}

// ControlForDirectWithEnhancedSecurity 控制模式 3[direct-with-enhanced-security]
func (c *Client) ControlForDirectWithEnhancedSecurity(objectRef string, value bool) error {
	return c.ControlByControlModel(objectRef, CONTROL_MODEL_DIRECT_ENHANCED, NewControlObjectParam(value))
}

// ControlForSboWithEnhancedSecurity 控制模式 4[sbo-with-enhanced-security]
func (c *Client) ControlForSboWithEnhancedSecurity(objectRef string, value bool) error {
	return c.ControlByControlModel(objectRef, CONTROL_MODEL_SBO_ENHANCED, NewControlObjectParam(value))
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
