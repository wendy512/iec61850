package scl

import "fmt"

type AttributeType int

const (
	Boolean AttributeType = iota
	Int8
	Int16
	Int32
	Int64
	Int128
	Int8U
	Int16U
	Int24U
	Int32U
	Float32
	Float64
	Enumerated
	OctetString64
	OctetString6
	OctetString8
	VisibleString32
	VisibleString64
	VisibleString65
	VisibleString129
	VisibleString255
	UnicodeString255
	Timestamp
	Quality
	Check
	CodedEnum
	GenericBitString
	Constructed
	EntryTime
	PhyComAddr
	Currency
	OptFlds
	TrgOps
)

func (that AttributeType) ToString() string {
	switch that {
	case Boolean:
		return "BOOLEAN"
	case Int8:
		return "INT8"
	case Int16:
		return "INT16"
	case Int32:
		return "INT32"
	case Int64:
		return "INT64"
	case Int128:
		return "INT128"
	case Int8U:
		return "INT8U"
	case Int16U:
		return "INT16U"
	case Int24U:
		return "INT24U"
	case Int32U:
		return "INT32U"
	case Float32:
		return "FLOAT32"
	case Float64:
		return "FLOAT64"
	case Enumerated:
		return "ENUMERATED"
	case OctetString64:
		return "OCTET_STRING_64"
	case OctetString6:
		return "OCTET_STRING_6"
	case OctetString8:
		return "OCTET_STRING_8"
	case VisibleString32:
		return "VISIBLE_STRING_32"
	case VisibleString64:
		return "VISIBLE_STRING_64"
	case VisibleString65:
		return "VISIBLE_STRING_65"
	case VisibleString129:
		return "VISIBLE_STRING_129"
	case VisibleString255:
		return "VISIBLE_STRING_255"
	case UnicodeString255:
		return "UNICODE_STRING_255"
	case Timestamp:
		return "TIMESTAMP"
	case Quality:
		return "QUALITY"
	case Check:
		return "CHECK"
	case CodedEnum:
		return "CODEDENUM"
	case GenericBitString:
		return "GENERIC_BITSTRING"
	case Constructed:
		return "CONSTRUCTED"
	case EntryTime:
		return "ENTRY_TIME"
	case PhyComAddr:
		return "PHYCOMADDR"
	case Currency:
		return "CURRENCY"
	case OptFlds:
		return "OPTFLDS"
	case TrgOps:
		return "TRGOPS"
	default:
		return "Unknown"
	}
}

func createAttributeTypeFromScl(typeString string) (AttributeType, error) {
	switch typeString {
	case "BOOLEAN":
		return Boolean, nil
	case "INT8":
		return Int8, nil
	case "INT16":
		return Int16, nil
	case "INT32":
		return Int32, nil
	case "INT64":
		return Int64, nil
	case "INT128":
		return Int128, nil
	case "INT8U":
		return Int8U, nil
	case "INT16U":
		return Int16U, nil
	case "INT24U":
		return Int24U, nil
	case "INT32U":
		return Int32U, nil
	case "FLOAT32":
		return Float32, nil
	case "FLOAT64":
		return Float64, nil
	case "Enum":
		return Enumerated, nil
	case "Dbpos":
		return CodedEnum, nil
	case "Check":
		return Check, nil
	case "Tcmd":
		return CodedEnum, nil
	case "Octet64":
		return OctetString64, nil
	case "Quality":
		return Quality, nil
	case "Timestamp":
		return Timestamp, nil
	case "Currency":
		return Currency, nil
	case "VisString32":
		return VisibleString32, nil
	case "VisString64":
		return VisibleString64, nil
	case "VisString65":
		return VisibleString65, nil
	case "VisString129":
		return VisibleString129, nil
	case "ObjRef":
		return VisibleString129, nil
	case "VisString255":
		return VisibleString255, nil
	case "Unicode255":
		return UnicodeString255, nil
	case "OptFlds":
		return OptFlds, nil
	case "TrgOps":
		return TrgOps, nil
	case "EntryID":
		return OctetString8, nil
	case "EntryTime":
		return EntryTime, nil
	case "PhyComAddr":
		return PhyComAddr, nil
	case "Struct":
		return Constructed, nil
	default:
		return -1, fmt.Errorf("unsupported attribute type %s", typeString)
	}
}

type SmpMod int

const (
	SmpPerPeriod SmpMod = iota
	SmpPerSecond
	SecPerSmp
)
