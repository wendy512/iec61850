//go:build linux && amd64

package iec61850

/*
#include "goose_publisher.h"
#include "mms_value.h"

static bool is_publisher_not_null(GoosePublisher p) {
	return p != NULL;
}

static void destroy_linked_list_val(LinkedList value) {
	LinkedList_destroyDeep(value, (LinkedListValueDeleteFunction)MmsValue_delete);
}
*/
import "C"
import (
	"errors"
	"unsafe"
)

type (
	LinkedListValue struct {
		internalLinkedList *C.struct_sLinkedList
	}

	GoosePublisherConf struct {
		InterfaceID  string
		AppID        uint16
		DstAddr      [6]uint8
		VlanID       uint16
		VlanPriority uint8
	}

	GoosePublisher struct {
		internalPublisher *C.struct_sGoosePublisher
	}
)

var (
	ErrCreateGoosePublisher = errors.New("can not create goose publisher")
	ErrSendGooseValue       = errors.New("can not send goose value")
)

func NewGoosePublisher(conf GoosePublisherConf) (publisher *GoosePublisher, err error) {
	parameters := C.struct_sCommParameters{}
	parameters.appId = C.uint16_t(conf.AppID)
	parameters.vlanId = C.uint16_t(conf.VlanID)
	parameters.vlanPriority = C.uint8_t(conf.VlanPriority)
	for i := 0; i < len(conf.DstAddr); i++ {
		parameters.dstAddress[i] = C.uint8_t(conf.DstAddr[i])
	}
	ether := C.CString(conf.InterfaceID)
	defer C.free(unsafe.Pointer(ether))

	cGoosePublisher := C.GoosePublisher_create(&parameters, ether)
	if !bool(C.is_publisher_not_null(cGoosePublisher)) {
		err = ErrCreateGoosePublisher
		return
	}

	publisher = &GoosePublisher{
		internalPublisher: cGoosePublisher,
	}

	return
}

func (receiver *GoosePublisher) SetGoCbRef(goCbRef string) {
	ref := C.CString(goCbRef)
	defer C.free(unsafe.Pointer(ref))

	C.GoosePublisher_setGoCbRef(receiver.internalPublisher, ref)
}

func (receiver *GoosePublisher) SetDataSetRef(dataSetRef string) {
	ref := C.CString(dataSetRef)
	defer C.free(unsafe.Pointer(ref))

	C.GoosePublisher_setDataSetRef(receiver.internalPublisher, ref)
}

func (receiver *GoosePublisher) SetConfRev(confRef uint32) {
	C.GoosePublisher_setConfRev(receiver.internalPublisher, C.uint32_t(confRef))
}

func (receiver *GoosePublisher) SetTimeAllowedToLive(timeAllowedToLive uint32) {
	C.GoosePublisher_setTimeAllowedToLive(receiver.internalPublisher, C.uint32_t(timeAllowedToLive))
}

func (receiver *GoosePublisher) SetSimulation(simulation bool) {
	C.GoosePublisher_setSimulation(receiver.internalPublisher, C.bool(simulation))
}

func (receiver *GoosePublisher) SetStNum(stNum uint32) {
	C.GoosePublisher_setStNum(receiver.internalPublisher, C.uint32_t(stNum))
}

func (receiver *GoosePublisher) SetSqNum(sqNum uint32) {
	C.GoosePublisher_setSqNum(receiver.internalPublisher, C.uint32_t(sqNum))
}

func (receiver *GoosePublisher) SetNeedsCommission(ndsCom bool) {
	C.GoosePublisher_setNeedsCommission(receiver.internalPublisher, C.bool(ndsCom))
}

func (receiver *GoosePublisher) IncreaseStNum() {
	C.GoosePublisher_increaseStNum(receiver.internalPublisher)
}

func (receiver *GoosePublisher) Reset() {
	C.GoosePublisher_reset(receiver.internalPublisher)
}

func (receiver *GoosePublisher) Publish(dataSet *LinkedListValue) error {
	if int(C.GoosePublisher_publish(receiver.internalPublisher, dataSet.internalLinkedList)) == -1 {
		return ErrSendGooseValue
	}

	return nil
}

func (receiver *GoosePublisher) Close() {
	C.GoosePublisher_destroy(receiver.internalPublisher)
}

func NewLinkedListValue() *LinkedListValue {
	return &LinkedListValue{
		internalLinkedList: C.LinkedList_create(),
	}
}

func (receiver *LinkedListValue) Add(value *MmsValue) error {
	rawVal, err := toMmsValue(value.Type, value.Value)
	if err != nil {
		return err
	}
	C.LinkedList_add(receiver.internalLinkedList, unsafe.Pointer(rawVal))

	return nil
}

func (receiver *LinkedListValue) Size() int {
	return int(C.LinkedList_size(receiver.internalLinkedList))
}

func (receiver *LinkedListValue) Destroy() {
	C.destroy_linked_list_val(receiver.internalLinkedList)
}
