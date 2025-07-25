package iec61850

import (
	"fmt"
	"time"
)

/*
#include "iec61850_common.h"
*/
import "C"

type (
	Timestamp struct {
		cTimestamp C.Timestamp
	}

	Quality  uint16
	Validity uint16
)

const (
	QUALITY_VALIDITY_GOOD         Quality = 0
	QUALITY_VALIDITY_INVALID      Quality = 2
	QUALITY_VALIDITY_RESERVED     Quality = 1
	QUALITY_VALIDITY_QUESTIONABLE Quality = 3
	QUALITY_DETAIL_OVERFLOW       Quality = 4
	QUALITY_DETAIL_OUT_OF_RANGE   Quality = 8
	QUALITY_DETAIL_BAD_REFERENCE  Quality = 16
	QUALITY_DETAIL_OSCILLATORY    Quality = 32
	QUALITY_DETAIL_FAILURE        Quality = 64
	QUALITY_DETAIL_OLD_DATA       Quality = 128
	QUALITY_DETAIL_INCONSISTENT   Quality = 256
	QUALITY_DETAIL_INACCURATE     Quality = 512
	QUALITY_SOURCE_SUBSTITUTED    Quality = 1024
	QUALITY_TEST                  Quality = 2048
	QUALITY_OPERATOR_BLOCKED      Quality = 4096
	QUALITY_DERIVED               Quality = 8192
)

const (
	VALIDITY_GOOD Validity = iota
	VALIDITY_INVALID
	VALIDITY_RESERVED
	VALIDITY_QUESTIONABLE
)

func (receiver Quality) GetValidity() Validity {
	return Validity(receiver & 0x3)
}

func NewTimestamp(time ...time.Time) *Timestamp {
	v := C.Timestamp_create()
	ret := &Timestamp{
		cTimestamp: *v,
	}
	C.Timestamp_destroy(v)
	switch len(time) {
	case 0:
		// skip
	case 1:
		ret.SetTime(time[0])
	default:
		panic(fmt.Errorf("expect got 0 or 1 time param, but got: %d", len(time)))
	}

	return ret
}

func (receiver *Timestamp) GetTimeInSeconds() uint32 {
	return uint32(C.Timestamp_getTimeInSeconds(&receiver.cTimestamp))
}

func (receiver *Timestamp) GetTimeInMs() uint64 {
	return uint64(C.Timestamp_getTimeInMs(&receiver.cTimestamp))
}

func (receiver *Timestamp) GetTimeInNs() uint64 {
	return uint64(C.Timestamp_getTimeInNs(&receiver.cTimestamp))
}

func (receiver *Timestamp) GetTime() time.Time {
	return time.Unix(0, int64(receiver.GetTimeInNs()))
}

func (receiver *Timestamp) IsLeapSecondKnown() bool {
	return bool(C.Timestamp_isLeapSecondKnown(&receiver.cTimestamp))
}

func (receiver *Timestamp) SetLeapSecondKnown(value bool) *Timestamp {
	C.Timestamp_setLeapSecondKnown(&receiver.cTimestamp, C.bool(value))
	return receiver
}

func (receiver *Timestamp) HasClockFailure() bool {
	return bool(C.Timestamp_hasClockFailure(&receiver.cTimestamp))
}

func (receiver *Timestamp) SetClockFailure(value bool) *Timestamp {
	C.Timestamp_setClockFailure(&receiver.cTimestamp, C.bool(value))

	return receiver
}

func (receiver *Timestamp) IsClockNotSynchronized() bool {
	return bool(C.Timestamp_isClockNotSynchronized(&receiver.cTimestamp))
}

func (receiver *Timestamp) SetClockNotSynchronized(value bool) *Timestamp {
	C.Timestamp_setClockNotSynchronized(&receiver.cTimestamp, C.bool(value))

	return receiver
}

func (receiver *Timestamp) GetSubSecondPrecision() int {
	return int(C.Timestamp_getSubsecondPrecision(&receiver.cTimestamp))
}

func (receiver *Timestamp) SetTime(time time.Time) *Timestamp {
	C.Timestamp_setTimeInNanoseconds(&receiver.cTimestamp, C.nsSinceEpoch(time.UnixNano()))
	return receiver
}
