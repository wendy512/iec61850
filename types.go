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

type CheckHandlerResult int

const (
	CONTROL_ACCEPTED                CheckHandlerResult = -1
	CONTROL_WAITING_FOR_SELECT      CheckHandlerResult = 0
	CONTROL_HARDWARE_FAULT          CheckHandlerResult = 1
	CONTROL_TEMPORARILY_UNAVAILABLE CheckHandlerResult = 2
	CONTROL_OBJECT_ACCESS_DENIED    CheckHandlerResult = 3
	CONTROL_OBJECT_UNDEFINED        CheckHandlerResult = 4
	CONTROL_VALUE_INVALID           CheckHandlerResult = 11
)

type ControlAddCause int

const (
	ADD_CAUSE_UNKNOWN                        ControlAddCause = 0
	ADD_CAUSE_NOT_SUPPORTED                  ControlAddCause = 1
	ADD_CAUSE_BLOCKED_BY_SWITCHING_HIERARCHY ControlAddCause = 2
	ADD_CAUSE_SELECT_FAILED                  ControlAddCause = 3
	ADD_CAUSE_INVALID_POSITION               ControlAddCause = 4
	ADD_CAUSE_POSITION_REACHED               ControlAddCause = 5
	ADD_CAUSE_PARAMETER_CHANGE_IN_EXECUTION  ControlAddCause = 6
	ADD_CAUSE_STEP_LIMIT                     ControlAddCause = 7
	ADD_CAUSE_BLOCKED_BY_MODE                ControlAddCause = 8
	ADD_CAUSE_BLOCKED_BY_PROCESS             ControlAddCause = 9
	ADD_CAUSE_BLOCKED_BY_INTERLOCKING        ControlAddCause = 10
	ADD_CAUSE_BLOCKED_BY_SYNCHROCHECK        ControlAddCause = 11
	ADD_CAUSE_COMMAND_ALREADY_IN_EXECUTION   ControlAddCause = 12
	ADD_CAUSE_BLOCKED_BY_HEALTH              ControlAddCause = 13
	ADD_CAUSE_1_OF_N_CONTROL                 ControlAddCause = 14
	ADD_CAUSE_ABORTION_BY_CANCEL             ControlAddCause = 15
	ADD_CAUSE_TIME_LIMIT_OVER                ControlAddCause = 16
	ADD_CAUSE_ABORTION_BY_TRIP               ControlAddCause = 17
	ADD_CAUSE_OBJECT_NOT_SELECTED            ControlAddCause = 18
	ADD_CAUSE_OBJECT_ALREADY_SELECTED        ControlAddCause = 19
	ADD_CAUSE_NO_ACCESS_AUTHORITY            ControlAddCause = 20
	ADD_CAUSE_ENDED_WITH_OVERSHOOT           ControlAddCause = 21
	ADD_CAUSE_ABORTION_DUE_TO_DEVIATION      ControlAddCause = 22
	ADD_CAUSE_ABORTION_BY_COMMUNICATION_LOSS ControlAddCause = 23
	ADD_CAUSE_ABORTION_BY_COMMAND            ControlAddCause = 24
	ADD_CAUSE_NONE                           ControlAddCause = 25
	ADD_CAUSE_INCONSISTENT_PARAMETERS        ControlAddCause = 26
	ADD_CAUSE_LOCKED_BY_OTHER_CLIENT         ControlAddCause = 27
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
