package iec61850

/*
#include "goose_receiver.h"
#include "goose_subscriber.h"
#include <stdint.h>

// 创建 goose 订阅对象
static GooseSubscriber create_simple_goose_subscriber(const char *goCbRef) {
	return GooseSubscriber_create(goCbRef, NULL);
}
*/
import "C"

import (
	"unsafe"
)

type (
	GooseReportCallback func(report *GooseReport)

	SubscriberConf struct {
		InterfaceID   string
		DstMacAddr    [6]uint8
		AppID         uint16
		Subscriber    string
		ReportHandler GooseReportCallback
	}

	GooseSubscriber struct {
		noCopy     struct{}
		Conf       SubscriberConf // 网卡名称
		subscriber *C.struct_sGooseSubscriber
		HandlerID  GooseCallbackHandlerID
	}

	GooseCallbackHandlerID uintptr
)

func NewGooseSubscriber(conf SubscriberConf) (subscriber *GooseSubscriber) {
	etherName := C.CString(conf.InterfaceID)
	defer C.free(unsafe.Pointer(etherName))
	goCbRef := C.CString(conf.Subscriber)
	defer C.free(unsafe.Pointer(goCbRef))

	cSubscriber := C.create_simple_goose_subscriber(goCbRef)
	C.GooseSubscriber_setDstMac(cSubscriber, (*C.uint8_t)(unsafe.Pointer(&conf.DstMacAddr[0])))
	C.GooseSubscriber_setAppId(cSubscriber, C.uint16_t(conf.AppID))
	newID := GooseCallbackHandlerID(gooseCallbackLocker.idOffset.Add(1))
	subscriber = &GooseSubscriber{
		subscriber: cSubscriber,
		Conf:       conf,
		HandlerID:  newID,
	}

	return
}

type (
	GooseReport struct {
		parameter unsafe.Pointer
		*GooseSubscriber
	}

	GooseParseError int
)

const (
	GooseParseErrorNoError GooseParseError = iota
	GooseParseErrorUnknownTag
	GooseParseErrorTagDecode
	GooseParseErrorSublevel
	GooseParseErrorOverflow
	GooseParseErrorUnderflow
	GooseParseErrorTypeMismatch
	GooseParseErrorLengthMismatch
	GooseParseErrorInvalidPadding
)

func (receiver *GooseSubscriber) GetGoID() string {
	return C.GoString(C.GooseSubscriber_getGoId(receiver.subscriber))
}

func (receiver *GooseSubscriber) GetGoCbRef() string {
	return C.GoString(C.GooseSubscriber_getGoCbRef(receiver.subscriber))
}

func (receiver *GooseSubscriber) GetDataSetName() string {
	return C.GoString(C.GooseSubscriber_getDataSet(receiver.subscriber))
}

func (receiver *GooseSubscriber) GetParseError() GooseParseError {
	return GooseParseError(C.GooseSubscriber_getParseError(receiver.subscriber))
}

func (receiver *GooseSubscriber) IsValid() bool {
	return bool(C.GooseSubscriber_isValid(receiver.subscriber))
}

func (receiver *GooseSubscriber) GetAppID() int32 {
	return int32(C.GooseSubscriber_getAppId(receiver.subscriber))
}

func (receiver *GooseSubscriber) GetSrcMac() [6]uint8 {
	bf := [6]byte{}
	C.GooseSubscriber_getSrcMac(receiver.subscriber, (*C.uint8_t)(unsafe.Pointer(&bf[0])))
	return bf
}

func (receiver *GooseSubscriber) GetDstMac() [6]uint8 {
	bf := [6]byte{}
	C.GooseSubscriber_getDstMac(receiver.subscriber, (*C.uint8_t)(unsafe.Pointer(&bf[0])))
	return bf
}

func (receiver *GooseSubscriber) GetStNum() uint32 {
	return uint32(C.GooseSubscriber_getStNum(receiver.subscriber))
}

func (receiver *GooseSubscriber) GetSqNum() uint32 {
	return uint32(C.GooseSubscriber_getSqNum(receiver.subscriber))
}

func (receiver *GooseSubscriber) IsTest() bool {
	return bool(C.GooseSubscriber_isTest(receiver.subscriber))
}

func (receiver *GooseSubscriber) GetConfRev() uint32 {
	return uint32(C.GooseSubscriber_getConfRev(receiver.subscriber))
}

func (receiver *GooseSubscriber) NeedsCommission() bool {
	return bool(C.GooseSubscriber_needsCommission(receiver.subscriber))
}

func (receiver *GooseSubscriber) GetTimeAllowedToLive() uint32 {
	return uint32(C.GooseSubscriber_getTimeAllowedToLive(receiver.subscriber))
}

func (receiver *GooseSubscriber) GetTimestamp() uint64 {
	return uint64(C.GooseSubscriber_getTimestamp(receiver.subscriber))
}

func (receiver *GooseSubscriber) IsVlanSet() bool {
	return bool(C.GooseSubscriber_isVlanSet(receiver.subscriber))
}

func (receiver *GooseSubscriber) GetVlanID() uint16 {
	return uint16(C.GooseSubscriber_getVlanId(receiver.subscriber))
}

func (receiver *GooseSubscriber) GetVlanPriority() uint8 {
	return uint8(C.GooseSubscriber_getVlanPrio(receiver.subscriber))
}

func (receiver *GooseSubscriber) GetDataSetValues() (*MmsValue, error) {
	cTypeMmsValue := C.GooseSubscriber_getDataSetValues(receiver.subscriber)
	mmsType := MmsType(C.MmsValue_getType(cTypeMmsValue))
	if mmsValue, err := toGoValue(cTypeMmsValue, mmsType); err != nil {
		return nil, err
	} else {
		return &MmsValue{
			Type:  mmsType,
			Value: mmsValue,
		}, nil
	}
}

func (receiver *GooseSubscriber) Destroy() {
	C.GooseSubscriber_destroy(receiver.subscriber)
	receiver.subscriber = nil
}
