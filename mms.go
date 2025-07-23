package iec61850

// #include <iec61850_client.h>
import "C"
import (
	"fmt"
	"unsafe"

	"github.com/spf13/cast"
)

func toMmsValue(mmsType MmsType, value interface{}) (*C.MmsValue, error) {
	var (
		mmsValue *C.MmsValue
		err      error
	)
	switch mmsType {
	case Boolean:
		mmsValue, err = toBoolMmsValue(value)
		if err != nil {
			return nil, err
		}
	case String:
		mmsValue, err = toStringMmsValue(value)
		if err != nil {
			return nil, err
		}
	case Float:
		mmsValue, err = toFloatMmsValue(value)
		if err != nil {
			return nil, err
		}
	case Uint8:
		mmsValue, err = toUint8MmsValue(value)
		if err != nil {
			return nil, err
		}
	case Uint16:
		mmsValue, err = toUint16MmsValue(value)
		if err != nil {
			return nil, err
		}
	case Uint32:
		mmsValue, err = toUint32MmsValue(value)
		if err != nil {
			return nil, err
		}
	case Int8:
		mmsValue, err = toInt8MmsValue(value)
		if err != nil {
			return nil, err
		}
	case Int16:
		mmsValue, err = toInt16MmsValue(value)
		if err != nil {
			return nil, err
		}
	case Int32:
		mmsValue, err = toInt32MmsValue(value)
		if err != nil {
			return nil, err
		}
	case Int64:
		mmsValue, err = toInt64MmsValue(value)
		if err != nil {
			return nil, err
		}
	default:
		return nil, UnSupportedOperation
	}
	return mmsValue, nil
}

func toGoValue(mmsValue *C.MmsValue, mmsType MmsType) (interface{}, error) {
	var (
		value interface{}
		err   error
	)

	switch mmsType {
	case Integer:
		value = int64(C.MmsValue_toInt64(mmsValue))
	case Unsigned:
		value = uint32(C.MmsValue_toUint32(mmsValue))
	case Boolean:
		value = bool(C.MmsValue_getBoolean(mmsValue))
	case Float:
		value = float32(C.MmsValue_toFloat(mmsValue))
	case String, VisibleString:
		value = C.GoString(C.MmsValue_toString(mmsValue))
	case Structure, Array:
		if value, err = toGoStructure(mmsValue, mmsType); err != nil {
			return nil, err
		}
	case BitString:
		value = uint32(C.MmsValue_getBitStringAsInteger(mmsValue))
	case OctetString:
		size := uint16(C.MmsValue_getOctetStringSize(mmsValue))
		bytes := make([]byte, size)
		for i := 0; i < int(size); i++ {
			bytes[i] = uint8(C.MmsValue_getOctetStringOctet(mmsValue, C.int(i)))
		}
		value = bytes
	case BinaryTime:
		value = uint64(C.MmsValue_getBinaryTimeAsUtcMs(mmsValue))
	case UTCTime:
		value = uint32(C.MmsValue_toUnixTimestamp(mmsValue))
	case DataAccessError:
		errorCode := C.MmsValue_getDataAccessError(mmsValue)
		return nil, fmt.Errorf("failed to read value (error code: %d)", int(errorCode))
	default:
		return nil, fmt.Errorf("unsupported type %d", mmsType)
	}
	return value, nil
}

func toGoStructure(mmsValue *C.MmsValue, mmsType MmsType) ([]*MmsValue, error) {
	if !(mmsType == Structure || mmsType == Array) {
		return nil, fmt.Errorf("require struct or array type value, but got type code is: %d", mmsType)
	}

	mmsValues := make([]*MmsValue, 0)
	for i := 0; ; i++ {
		value := C.MmsValue_getElement(mmsValue, C.int(i))
		// 读不到表示节点下没有属性了
		if value == nil {
			return mmsValues, nil
		}

		valueType := MmsType(C.MmsValue_getType(value))
		goValue, err := toGoValue(value, valueType)
		if err != nil {
			return nil, err
		}

		mmsValues = append(mmsValues, &MmsValue{
			Type:  valueType,
			Value: goValue,
		})
	}
}

func toInt64MmsValue(value interface{}) (*C.MmsValue, error) {
	v, err := cast.ToInt64E(value)
	if err != nil {
		return nil, err
	}
	// int64
	return C.MmsValue_newIntegerFromInt64(C.int64_t(v)), nil
}

func toInt32MmsValue(value interface{}) (*C.MmsValue, error) {
	v, err := cast.ToInt32E(value)
	if err != nil {
		return nil, err
	}
	// int32
	return C.MmsValue_newIntegerFromInt32(C.int32_t(v)), nil
}

func toInt16MmsValue(value interface{}) (*C.MmsValue, error) {
	v, err := cast.ToInt16E(value)
	if err != nil {
		return nil, err
	}
	// int16
	return C.MmsValue_newIntegerFromInt16(C.int16_t(v)), nil
}

func toInt8MmsValue(value interface{}) (*C.MmsValue, error) {
	v, err := cast.ToInt8E(value)
	if err != nil {
		return nil, err
	}
	// int8
	return C.MmsValue_newIntegerFromInt8(C.int8_t(v)), nil
}

func toUint32MmsValue(value interface{}) (*C.MmsValue, error) {
	v, err := cast.ToUint32E(value)
	if err != nil {
		return nil, err
	}
	// uint32
	mmsValue := C.MmsValue_newUnsigned(C.int(32))
	C.MmsValue_setUint32(mmsValue, C.uint32_t(v))
	return mmsValue, nil
}

func toUint16MmsValue(value interface{}) (*C.MmsValue, error) {
	v, err := cast.ToUint16E(value)
	if err != nil {
		return nil, err
	}
	// uint16
	mmsValue := C.MmsValue_newUnsigned(C.int(16))
	C.MmsValue_setUint16(mmsValue, C.uint16_t(v))
	return mmsValue, nil
}

func toUint8MmsValue(value interface{}) (*C.MmsValue, error) {
	v, err := cast.ToUint8E(value)
	if err != nil {
		return nil, err
	}
	// uint8
	mmsValue := C.MmsValue_newUnsigned(C.int(8))
	C.MmsValue_setUint8(mmsValue, C.uint8_t(v))
	return mmsValue, nil
}

func toFloatMmsValue(value interface{}) (*C.MmsValue, error) {
	v, err := cast.ToFloat32E(value)
	if err != nil {
		return nil, err
	}
	return C.MmsValue_newFloat(C.float(v)), nil
}

func toBoolMmsValue(value interface{}) (*C.MmsValue, error) {
	v, err := cast.ToBoolE(value)
	if err != nil {
		return nil, err
	}
	return C.MmsValue_newBoolean(C.bool(v)), nil
}

func toStringMmsValue(value interface{}) (*C.MmsValue, error) {
	v, err := cast.ToStringE(value)
	if err != nil {
		return nil, err
	}
	stringValue := C.CString(v)
	defer C.free(unsafe.Pointer(stringValue))
	mmsValue := C.MmsValue_newMmsString(stringValue)
	return mmsValue, nil
}
