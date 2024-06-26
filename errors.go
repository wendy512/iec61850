package iec61850

// #include <iec61850_client.h>
import "C"
import "errors"

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
	StructureMustBeMmsValue           = errors.New("structure type must be MmsValue array")
	CreateControlObjectClientFail     = errors.New("control object not found in server")
	ControlObjectFail                 = errors.New("control object fail")
	ControlSelectFail                 = errors.New("select control fail")
	UnSupportedOperation              = errors.New("unsupported operation")
)

func GetIedClientError(err C.IedClientError) error {
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
