package iec61850

/*
#cgo CFLAGS: -I${SRCDIR}/libiec61850/head
#cgo LDFLAGS: -L${SRCDIR}/libiec61850/win64 -liec61850 -lws2_32
#include <stdlib.h>
#include <stdio.h>
#include "iec61850_client.h"
*/
import "C"
import (
	"errors"
	"github.com/spf13/cast"
	"sync/atomic"
	"unsafe"
)

// fc types
const (
	// ST Status information
	ST FC = "ST"
	// MX Measurands - analogue values
	MX FC = "MX"
	// SP Setpoint
	SP FC = "SP"
	// SV Substitution
	SV FC = "SV"
	// CF Configuration
	CF FC = "CF"
	// DC Description
	DC FC = "DC"
	// SG Setting group
	SG FC = "SG"
	// SE Setting group editable
	SE FC = "SE"
	// SR service response / service tracking
	SR FC = "SR"
	// OR Operate received
	OR FC = "OR"
	// BL Blocking
	BL FC = "BL"
	// EX Extended definition
	EX FC = "EX"
	// CO Control, deprecated but kept here for backward compatibility
	CO FC = "CO"
	// RP Unbuffered Reporting
	RP FC = "RP"
	// BR Buffered Reporting
	BR FC = "BR"
)

// data types
const (
	Array Type = iota
	Structure
	Boolean
	BitString
	Int8
	Int16
	Int32
	Int64
	Uint8
	Uint16
	Uint32
	Float
	OctetString
	VisibleString
	GeneralizedTime
	BinaryTime
	Bcd
	ObjId
	String
	UTCTime
	DataAccessError
)

var (
	NotConnected                      = errors.New("the service request can not be executed because the client is not yet connected")
	AlreadyConnected                  = errors.New("connect service not execute because the client is already connected")
	ConnectionLost                    = errors.New("the service request can not be executed caused by a loss of connection")
	ServiceNotSupported               = errors.New("the service or some given parameters are not supported by the client stack or by the server")
	ConnectionRejected                = errors.New("connection rejected by server")
	OutstandingCallLimitReached       = errors.New("cannot send request because outstanding call limit is reached")
	UserProvidedInvalidArgument       = errors.New("API function has been called with an invalid argument")
	EnableReportFailedDatasetMismatch = errors.New("API function has been called with an invalid argument")
	ObjectReferenceInvalid            = errors.New("the object provided object reference is invalid (there is a syntactical error)")
	UnexpectedValueReceived           = errors.New("received object is of unexpected type")
	Timeout                           = errors.New("the communication to the server failed with a timeout")
	AccessDenied                      = errors.New("the server rejected the access to the requested object/service due to access control")
	ObjectDoesNotExist                = errors.New("the server reported that the requested object does not exist (returned by server)")
	ObjectExists                      = errors.New("the server reported that the requested object already exists")
	ObjectAccessUnsupported           = errors.New("the server does not support the requested access method (returned by server)")
	TypeInconsistent                  = errors.New("the server expected an object of another type (returned by server)")
	TemporarilyUnavailable            = errors.New("the object or service is temporarily unavailable (returned by server)")
	ObjectUndefined                   = errors.New("the specified object is not defined in the server (returned by server)")
	InvalidAddress                    = errors.New("the specified address is invalid (returned by server)")
	HardwareFault                     = errors.New("service failed due to a hardware fault (returned by server)")
	TypeUnsupported                   = errors.New("the requested data type is not supported by the server (returned by server)")
	ObjectAttributeInconsistent       = errors.New("the provided attributes are inconsistent (returned by server)")
	ObjectValueInvalid                = errors.New("the provided object value is invalid (returned by server)")
	ObjectInvalidated                 = errors.New("the object is invalidated (returned by server)")
	MalformedMessage                  = errors.New("received an invalid response message from the server")
	ServiceNotImplemented             = errors.New("service not implemented")
	Unknown                           = errors.New("unknown error")
	UnSupportOperation                = errors.New("un support operation")
)

type Type uint

type FC string

type Connection struct {
	conn      C.IedConnection
	connected *atomic.Bool
}

type Settings struct {
	Host           string
	Port           int
	ConnectTimeout uint
	RequestTimeout uint
}

func NewSettings() *Settings {
	return &Settings{
		Host:           "localhost",
		Port:           102,
		ConnectTimeout: 10000,
		RequestTimeout: 10000,
	}
}

func NewConnection(settings *Settings) (*Connection, error) {
	conn, clientErr := connect(settings)
	if err := getIedClientError(clientErr); err != nil {
		return nil, err
	}

	connected := &atomic.Bool{}
	connected.Store(true)
	connection := &Connection{
		conn:      conn,
		connected: connected,
	}
	return connection, nil
}

func (c *Connection) Write(objectReference string, fc FC, value interface{}) error {
	if !c.connected.Load() {
		return NotConnected
	}

	specType := getVariableSpecType(c.conn, objectReference, fc)
	var (
		mmsValue    *C.MmsValue
		clientError C.IedClientError
	)

	switch specType {
	case Boolean:
		v, err := cast.ToBoolE(value)
		if err != nil {
			return err
		}
		mmsValue = C.MmsValue_newBoolean(C.bool(cast.ToInt(v)))
	case String:
		v, err := cast.ToStringE(value)
		if err != nil {
			return err
		}
		stringValue := C.CString(v)
		// 释放内存
		defer C.free(unsafe.Pointer(stringValue))
		mmsValue = C.MmsValue_newMmsString(stringValue)
	case Float:
		v, err := cast.ToFloat32E(value)
		if err != nil {
			return err
		}
		mmsValue = C.MmsValue_newFloat(C.float(v))
	case Uint8:
		v, err := cast.ToUint8E(value)
		if err != nil {
			return err
		}
		// uint8
		mmsValue = C.MmsValue_newUnsigned(C.int(1))
		C.MmsValue_setUint8(mmsValue, C.uint8_t(v))
	case Uint16:
		v, err := cast.ToUint16E(value)
		if err != nil {
			return err
		}
		// uint16
		mmsValue = C.MmsValue_newUnsigned(C.int(2))
		C.MmsValue_setUint16(mmsValue, C.uint16_t(v))
	case Uint32:
		v, err := cast.ToUint32E(value)
		if err != nil {
			return err
		}
		// uint32
		mmsValue = C.MmsValue_newUnsigned(C.int(4))
		C.MmsValue_setUint32(mmsValue, C.uint32_t(v))
	case Int8:
		v, err := cast.ToInt8E(value)
		if err != nil {
			return err
		}
		// int8
		mmsValue = C.MmsValue_newIntegerFromInt8(C.int8_t(v))
	case Int16:
		v, err := cast.ToInt16E(value)
		if err != nil {
			return err
		}
		// int16
		mmsValue = C.MmsValue_newIntegerFromInt16(C.int16_t(v))
	case Int32:
		v, err := cast.ToInt32E(value)
		if err != nil {
			return err
		}
		// int32
		mmsValue = C.MmsValue_newIntegerFromInt32(C.int32_t(v))
	case Int64:
		v, err := cast.ToInt64E(value)
		if err != nil {
			return err
		}
		// int64
		mmsValue = C.MmsValue_newIntegerFromInt64(C.int64_t(v))
	default:
		return UnSupportOperation
	}

	cObjectReference := C.CString(objectReference)
	defer C.free(unsafe.Pointer(cObjectReference))
	C.IedConnection_writeObject(c.conn, &clientError, cObjectReference, fc.getFunctionalConstraint(), mmsValue)
	C.MmsValue_delete(mmsValue)
	return getIedClientError(clientError)
}

func (c *Connection) Read(objectReference string, fc FC) (interface{}, error) {
	var (
		clientError C.IedClientError
		value       interface{}
	)
	cObjectReference := C.CString(objectReference)
	defer C.free(unsafe.Pointer(cObjectReference))

	getVariableSpecType(c.conn, objectReference, fc)
	mmsValue := C.IedConnection_readObject(c.conn, &clientError, cObjectReference, fc.getFunctionalConstraint())
	if err := getIedClientError(clientError); err != nil {
		return nil, err
	}

	mmsType := C.MmsValue_getType(mmsValue)
	switch mmsType {
	case C.MMS_FLOAT:
		value = float32(C.MmsValue_toFloat(mmsValue))
	case C.MMS_BOOLEAN:
		value = int(C.MmsValue_toInt32(mmsValue))
	case C.MMS_INTEGER:
		value = int(C.MmsValue_toInt32(mmsValue))
	case C.MMS_UNSIGNED:
		value = uint(C.MmsValue_toUint32(mmsValue))
	}
	C.MmsValue_delete(mmsValue)
	return value, nil
}

func (c *Connection) Close() {
	if c.conn != nil && c.connected.CompareAndSwap(true, false) {
		C.IedConnection_destroy(c.conn)
	}
}

func getIedClientError(err C.IedClientError) error {
	cError := C.IedClientError(err)
	switch cError {
	case C.IED_ERROR_OK:
		return nil
	case C.IED_ERROR_NOT_CONNECTED:
		return NotConnected
	case C.IED_ERROR_ALREADY_CONNECTED:
		return AlreadyConnected
	case C.IED_ERROR_CONNECTION_LOST:
		return ConnectionLost
	case C.IED_ERROR_SERVICE_NOT_SUPPORTED:
		return ServiceNotSupported
	case C.IED_ERROR_CONNECTION_REJECTED:
		return ConnectionRejected
	case C.IED_ERROR_OUTSTANDING_CALL_LIMIT_REACHED:
		return OutstandingCallLimitReached
	case C.IED_ERROR_USER_PROVIDED_INVALID_ARGUMENT:
		return UserProvidedInvalidArgument
	case C.IED_ERROR_ENABLE_REPORT_FAILED_DATASET_MISMATCH:
		return EnableReportFailedDatasetMismatch
	case C.IED_ERROR_OBJECT_REFERENCE_INVALID:
		return ObjectReferenceInvalid
	case C.IED_ERROR_UNEXPECTED_VALUE_RECEIVED:
		return UnexpectedValueReceived
	case C.IED_ERROR_TIMEOUT:
		return Timeout
	case C.IED_ERROR_ACCESS_DENIED:
		return AccessDenied
	case C.IED_ERROR_OBJECT_DOES_NOT_EXIST:
		return ObjectDoesNotExist
	case C.IED_ERROR_OBJECT_EXISTS:
		return ObjectExists
	case C.IED_ERROR_OBJECT_ACCESS_UNSUPPORTED:
		return ObjectAccessUnsupported
	case C.IED_ERROR_TYPE_INCONSISTENT:
		return TypeInconsistent
	case C.IED_ERROR_TEMPORARILY_UNAVAILABLE:
		return TemporarilyUnavailable
	case C.IED_ERROR_OBJECT_UNDEFINED:
		return ObjectUndefined
	case C.IED_ERROR_INVALID_ADDRESS:
		return InvalidAddress
	case C.IED_ERROR_HARDWARE_FAULT:
		return HardwareFault
	case C.IED_ERROR_TYPE_UNSUPPORTED:
		return TypeUnsupported
	case C.IED_ERROR_OBJECT_ATTRIBUTE_INCONSISTENT:
		return ObjectAttributeInconsistent
	case C.IED_ERROR_OBJECT_VALUE_INVALID:
		return ObjectValueInvalid
	case C.IED_ERROR_OBJECT_INVALIDATED:
		return ObjectInvalidated
	case C.IED_ERROR_MALFORMED_MESSAGE:
		return MalformedMessage
	case C.IED_ERROR_SERVICE_NOT_IMPLEMENTED:
		return ServiceNotImplemented
	default:
		return Unknown
	}
}

func (fc FC) getFunctionalConstraint() C.FunctionalConstraint {
	var cfc C.FunctionalConstraint = C.IEC61850_FC_ST
	switch fc {
	case ST:
		cfc = C.IEC61850_FC_ST
	case MX:
		cfc = C.IEC61850_FC_MX
	case SP:
		cfc = C.IEC61850_FC_SP
	case SV:
		cfc = C.IEC61850_FC_SV
	case CF:
		cfc = C.IEC61850_FC_CF
	case DC:
		cfc = C.IEC61850_FC_DC
	case SG:
		cfc = C.IEC61850_FC_SG
	case SE:
		cfc = C.IEC61850_FC_SE
	case SR:
		cfc = C.IEC61850_FC_SR
	case OR:
		cfc = C.IEC61850_FC_OR
	case BL:
		cfc = C.IEC61850_FC_BL
	case EX:
		cfc = C.IEC61850_FC_EX
	case CO:
		cfc = C.IEC61850_FC_CO
	case RP:
		cfc = C.IEC61850_FC_RP
	case BR:
		cfc = C.IEC61850_FC_BR
	}
	return cfc
}

func getVariableSpecType(conn C.IedConnection, objectReference string, fc FC) Type {
	var err C.IedClientError
	cObjectReference := C.CString(objectReference)
	defer C.free(unsafe.Pointer(cObjectReference))
	// 获取类型
	spec := C.IedConnection_getVariableSpecification(conn, &err, cObjectReference, fc.getFunctionalConstraint())
	mssType := C.MmsVariableSpecification_getType(spec)
	switch mssType {
	case C.MMS_ARRAY:
		return Array
	case C.MMS_STRUCTURE:
		return Structure
	case C.MMS_BOOLEAN:
		return Boolean
	case C.MMS_BIT_STRING:
		return BitString
	case C.MMS_INTEGER:
		i := int(spec.typeSpec[3])
		switch i {
		case 1:
			return Int8
		case 2:
			return Int16
		case 4:
			return Int32
		default:
			return Int64
		}
	case C.MMS_UNSIGNED:
		switch int(spec.typeSpec[4]) {
		case 1:
			return Uint8
		case 2:
			return Uint16
		default:
			return Uint32
		}
	case C.MMS_FLOAT:
		return Float
	case C.MMS_OCTET_STRING:
		return OctetString
	case C.MMS_VISIBLE_STRING:
		return VisibleString
	case C.MMS_GENERALIZED_TIME:
		return GeneralizedTime
	case C.MMS_BINARY_TIME:
		return BinaryTime
	case C.MMS_BCD:
		return Bcd
	case C.MMS_OBJ_ID:
		return ObjId
	case C.MMS_STRING:
		return String
	case C.MMS_UTC_TIME:
		return UTCTime
	default:
		return DataAccessError
	}
}

func connect(settings *Settings) (C.IedConnection, C.IedClientError) {
	conn := C.IedConnection_create()
	C.IedConnection_setConnectTimeout(conn, C.uint(settings.ConnectTimeout))
	C.IedConnection_setRequestTimeout(conn, C.uint(settings.RequestTimeout))
	host := C.CString(settings.Host)
	// 释放内存
	defer C.free(unsafe.Pointer(host))
	var err C.IedClientError = C.IED_ERROR_OK
	// 建立连接
	C.IedConnection_connect(conn, &err, host, C.int(settings.Port))
	return conn, err
}
