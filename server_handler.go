package iec61850

/*
#include <iec61850_server.h>

typedef struct {
    uint8_t* buf;
    int length;
} Buffer;

extern MmsDataAccessError writeAccessHandlerBridge(DataAttribute* dataAttribute, MmsValue* value, ClientConnection connection, void* parameter);

extern ControlHandlerResult controlHandlerBridge(ControlAction action, void* parameter, MmsValue* ctlVal, bool test);

extern bool acseAuthenticatorBridge(void* parameter, AcseAuthenticationParameter authParameter, void** securityToken, IsoApplicationReference* appReference);

static Buffer AcseAuthenticationParameter_GetBuffer(AcseAuthenticationParameter authParameter) {
    if (authParameter->mechanism == ACSE_AUTH_PASSWORD) {
        uint8_t *buf = authParameter->value.password.octetString;
		int len = authParameter->value.password.passwordLength;
		return (Buffer){buf, len};
    } else if (authParameter->mechanism == ACSE_AUTH_CERTIFICATE || authParameter->mechanism == ACSE_AUTH_TLS) {
 		uint8_t *buf = authParameter->value.certificate.buf;
		int len = authParameter->value.certificate.length;
		return (Buffer){buf, len};
    }
    return (Buffer){NULL, 0};
}
*/
import "C"

import (
	"fmt"
	"sync/atomic"
	"unsafe"
)

var (
	callbackIdGen        = atomic.Int32{}
	writeAccessCallbacks = make(map[int32]*writeAccessCallback)
	controlCallbacks     = make(map[int32]*controlCallback)
)

type writeAccessCallback struct {
	node    *ModelNode
	handler WriteAccessHandler
}

type controlCallback struct {
	node    *ModelNode
	handler ControlHandler
}

type ControlAction struct {
	ControlTime    uint64
	IsSelect       bool
	InterlockCheck bool
	SynchroCheck   bool
	CtlNum         int
	OrIdent        []byte
	OrCat          int
}

type IsoApplicationReference struct {
	ApTitle     []uint16
	AeQualifier int
}

type AcseAuthenticationParameter struct {
	Mechanism   AcseAuthenticationMechanism
	Password    []byte // for mechanism = ACSE_AUTH_PASSWORD
	Certificate []byte // for mechanism = ACSE_AUTH_CERTIFICATE or ACSE_AUTH_TLS
}

type WriteAccessHandler func(node *ModelNode, mmsValue *MmsValue) MmsDataAccessError

type ControlHandler func(node *ModelNode, action *ControlAction, mmsValue *MmsValue, test bool) ControlHandlerResult

type ClientAuthenticator func(securityToken *unsafe.Pointer, authParameter *AcseAuthenticationParameter, appReference *IsoApplicationReference) bool

//export writeAccessHandlerBridge
func writeAccessHandlerBridge(dataAttribute *C.DataAttribute, value *C.MmsValue, connection C.ClientConnection, parameter unsafe.Pointer) C.MmsDataAccessError {
	callbackId := int32(uintptr(parameter))
	if call, ok := writeAccessCallbacks[callbackId]; ok {

		mmsType := MmsType(C.MmsValue_getType(value))
		if goValue, err := toGoValue(value, mmsType); err == nil {

			dataAccessError := call.handler(call.node, &MmsValue{
				Type:  mmsType,
				Value: goValue,
			})
			return C.MmsDataAccessError(dataAccessError)
		} else {
			fmt.Printf("mms value to go value error: %v\n", err)
		}
	}
	return C.DATA_ACCESS_ERROR_OBJECT_ACCESS_DENIED
}

//export controlHandlerBridge
func controlHandlerBridge(action C.ControlAction, parameter unsafe.Pointer, ctlVal *C.MmsValue, test C.bool) C.ControlHandlerResult {
	callbackId := int32(uintptr(parameter))
	if call, ok := controlCallbacks[callbackId]; ok {

		mmsType := MmsType(C.MmsValue_getType(ctlVal))
		if goValue, err := toGoValue(ctlVal, mmsType); err == nil {

			var (
				orIdentSize C.int
				orIdent     []byte
			)

			orIdentBuffer := C.ControlAction_getOrIdent(action, (*C.int)(unsafe.Pointer(&orIdentSize)))
			if orIdentBuffer != nil {
				size := int(orIdentSize)
				orIdent = C.GoBytes(unsafe.Pointer(orIdentBuffer), C.int(size))
			}

			actionFill := &ControlAction{
				ControlTime:    uint64(C.ControlAction_getControlTime(action)),
				IsSelect:       bool(C.ControlAction_isSelect(action)),
				InterlockCheck: bool(C.ControlAction_getInterlockCheck(action)),
				SynchroCheck:   bool(C.ControlAction_getSynchroCheck(action)),
				CtlNum:         int(C.ControlAction_getCtlNum(action)),
				OrIdent:        orIdent,
				OrCat:          int(C.ControlAction_getOrCat(action)),
			}

			controlHandlerResult := call.handler(call.node, actionFill, &MmsValue{mmsType, goValue}, bool(test))
			return C.ControlHandlerResult(controlHandlerResult)
		}
	}
	return C.CONTROL_RESULT_FAILED
}

//export acseAuthenticatorBridge
func acseAuthenticatorBridge(parameter unsafe.Pointer, authParameter C.AcseAuthenticationParameter, securityToken *unsafe.Pointer, appReference *C.IsoApplicationReference) C.bool {
	is := (*IedServer)(parameter)

	oid := appReference.apTitle
	// 将 oid->arc 转换为 Go 的切片
	arcValues := make([]uint16, int(oid.arcCount))
	for i := 0; i < int(oid.arcCount); i++ {
		arcValues[i] = uint16(oid.arc[i])
	}

	_appReference := &IsoApplicationReference{
		ApTitle:     arcValues,
		AeQualifier: int(appReference.aeQualifier),
	}

	mechanism := AcseAuthenticationMechanism(int(authParameter.mechanism))
	_authParameter := &AcseAuthenticationParameter{
		Mechanism: mechanism,
	}

	// 将 uint8_t* 转换为 Go 的 []byte
	buffer := C.AcseAuthenticationParameter_GetBuffer(authParameter)
	switch mechanism {
	case ACSE_AUTH_PASSWORD:
		_authParameter.Password = C.GoBytes(unsafe.Pointer(buffer.buf), C.int(buffer.length))
	case ACSE_AUTH_CERTIFICATE, ACSE_AUTH_TLS:
		_authParameter.Certificate = C.GoBytes(unsafe.Pointer(buffer.buf), C.int(buffer.length))
	default:
		// none
	}

	result := is.clientAuthenticator(securityToken, _authParameter, _appReference)
	return C.bool(result)
}

func (is *IedServer) SetHandleWriteAccess(modelNode *ModelNode, handler WriteAccessHandler) {
	if modelNode == nil {
		return
	}

	callbackId := callbackIdGen.Add(1)
	// 将 int 转为 uintptr，再转为 unsafe.Pointer
	cPtr := intToPointerBug58625(callbackId)
	writeAccessCallbacks[callbackId] = &writeAccessCallback{
		node:    modelNode,
		handler: handler,
	}

	C.IedServer_handleWriteAccess(is.server, (*C.DataAttribute)(modelNode._modelNode), (*[0]byte)(C.writeAccessHandlerBridge), cPtr)
}

func (is *IedServer) SetControlHandler(modelNode *ModelNode, handler ControlHandler) {
	if modelNode == nil {
		return
	}

	callbackId := callbackIdGen.Add(1)
	// 将 int 转为 uintptr，再转为 unsafe.Pointer
	cPtr := intToPointerBug58625(callbackId)
	controlCallbacks[callbackId] = &controlCallback{
		node:    modelNode,
		handler: handler,
	}

	C.IedServer_setControlHandler(is.server, (*C.DataObject)(modelNode._modelNode), (*[0]byte)(C.controlHandlerBridge), cPtr)
}

// intToPointerBug58625 is a helper function to fix issue #58625 in Go | https://github.com/golang/go/issues/58625
func intToPointerBug58625(i int32) unsafe.Pointer {
	var intPtr = uintptr(i)
	return *(*unsafe.Pointer)(unsafe.Pointer(&intPtr))
}

func (is *IedServer) SetAuthenticator(clientAuthenticator ClientAuthenticator) {
	is.clientAuthenticator = clientAuthenticator
	cPtr := unsafe.Pointer(is)
	C.IedServer_setAuthenticator(is.server, (*[0]byte)(C.acseAuthenticatorBridge), cPtr)
}
