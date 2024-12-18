package iec61850

type MmsType int

type MmsValue struct {
	Type  MmsType
	Value interface{}
}

// data types
const (
	Array MmsType = iota
	Structure
	Boolean
	BitString
	Integer
	Unsigned
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
	Int8
	Int16
	Int32
	Int64
	Uint8
	Uint16
	Uint32
)

type MmsDataAccessError int

const (
	DATA_ACCESS_ERROR_SUCCESS_NO_UPDATE             MmsDataAccessError = -3
	DATA_ACCESS_ERROR_NO_RESPONSE                                      = -2
	DATA_ACCESS_ERROR_SUCCESS                                          = -1
	DATA_ACCESS_ERROR_OBJECT_INVALIDATED                               = 0
	DATA_ACCESS_ERROR_HARDWARE_FAULT                                   = 1
	DATA_ACCESS_ERROR_TEMPORARILY_UNAVAILABLE                          = 2
	DATA_ACCESS_ERROR_OBJECT_ACCESS_DENIED                             = 3
	DATA_ACCESS_ERROR_OBJECT_UNDEFINED                                 = 4
	DATA_ACCESS_ERROR_INVALID_ADDRESS                                  = 5
	DATA_ACCESS_ERROR_TYPE_UNSUPPORTED                                 = 6
	DATA_ACCESS_ERROR_TYPE_INCONSISTENT                                = 7
	DATA_ACCESS_ERROR_OBJECT_ATTRIBUTE_INCONSISTENT                    = 8
	DATA_ACCESS_ERROR_OBJECT_ACCESS_UNSUPPORTED                        = 9
	DATA_ACCESS_ERROR_OBJECT_NONE_EXISTENT                             = 10
	DATA_ACCESS_ERROR_OBJECT_VALUE_INVALID                             = 11
	DATA_ACCESS_ERROR_UNKNOWN                                          = 12
)

type ControlHandlerResult int

const (
	CONTROL_RESULT_FAILED ControlHandlerResult = iota
	CONTROL_RESULT_OK
	CONTROL_RESULT_WAITING
)

type ControlModel int

const (
	// CONTROL_MODEL_STATUS_ONLY No support for control functions. Control object only support status information.
	CONTROL_MODEL_STATUS_ONLY ControlModel = iota
	// CONTROL_MODEL_DIRECT_NORMAL Direct control with normal security: Supports Operate, TimeActivatedOperate (optional), and Cancel (optional).
	CONTROL_MODEL_DIRECT_NORMAL
	// CONTROL_MODEL_SBO_NORMAL Select before operate (SBO) with normal security: Supports Select, Operate, TimeActivatedOperate (optional), and Cancel (optional).
	CONTROL_MODEL_SBO_NORMAL
	// CONTROL_MODEL_DIRECT_ENHANCED Direct control with enhanced security (enhanced security includes the CommandTermination service)
	CONTROL_MODEL_DIRECT_ENHANCED
	// CONTROL_MODEL_SBO_ENHANCED Select before operate (SBO) with enhanced security (enhanced security includes the CommandTermination service)
	CONTROL_MODEL_SBO_ENHANCED
)

type AcseAuthenticationMechanism int

const (
	// ACSE_AUTH_NONE Neither ACSE nor TLS authentication used
	ACSE_AUTH_NONE AcseAuthenticationMechanism = iota

	// ACSE_AUTH_PASSWORD Use ACSE password for client authentication
	ACSE_AUTH_PASSWORD

	// ACSE_AUTH_CERTIFICATE Use ACSE certificate for client authentication
	ACSE_AUTH_CERTIFICATE

	// ACSE_AUTH_TLS Use TLS certificate for client authentication
	ACSE_AUTH_TLS
)
