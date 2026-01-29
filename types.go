package iec61850

type MmsType int

type MmsValue struct {
	Type  MmsType
	Value interface{}
}

// UtcTimeValue holds a UTC time from an MMS UTCTime with millisecond precision and time quality.
// Returned when reading a UTCTime attribute via Client.Read() or related APIs.
type UtcTimeValue struct {
	Milliseconds uint64 // Milliseconds since Unix epoch (1970-01-01 00:00:00 UTC)
	TimeQuality  uint8  // IEC 61850 time quality (leapSecondsKnown, clockFailure, clockNotSynchronized, subsecond accuracy)
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
	DATA_ACCESS_ERROR_NO_RESPONSE                   MmsDataAccessError = -2
	DATA_ACCESS_ERROR_SUCCESS                       MmsDataAccessError = -1
	DATA_ACCESS_ERROR_OBJECT_INVALIDATED            MmsDataAccessError = 0
	DATA_ACCESS_ERROR_HARDWARE_FAULT                MmsDataAccessError = 1
	DATA_ACCESS_ERROR_TEMPORARILY_UNAVAILABLE       MmsDataAccessError = 2
	DATA_ACCESS_ERROR_OBJECT_ACCESS_DENIED          MmsDataAccessError = 3
	DATA_ACCESS_ERROR_OBJECT_UNDEFINED              MmsDataAccessError = 4
	DATA_ACCESS_ERROR_INVALID_ADDRESS               MmsDataAccessError = 5
	DATA_ACCESS_ERROR_TYPE_UNSUPPORTED              MmsDataAccessError = 6
	DATA_ACCESS_ERROR_TYPE_INCONSISTENT             MmsDataAccessError = 7
	DATA_ACCESS_ERROR_OBJECT_ATTRIBUTE_INCONSISTENT MmsDataAccessError = 8
	DATA_ACCESS_ERROR_OBJECT_ACCESS_UNSUPPORTED     MmsDataAccessError = 9
	DATA_ACCESS_ERROR_OBJECT_NONE_EXISTENT          MmsDataAccessError = 10
	DATA_ACCESS_ERROR_OBJECT_VALUE_INVALID          MmsDataAccessError = 11
	DATA_ACCESS_ERROR_UNKNOWN                       MmsDataAccessError = 12
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
