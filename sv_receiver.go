//go:build linux && amd64

package iec61850

/*
#include "sv_subscriber.h"
#include <stdint.h>

extern void svCGOReportHandler(void* parameter, SVSubscriber_ASDU asdu);

static void svListenerProxy(SVSubscriber subscriber, void* parameter, SVSubscriber_ASDU asdu) {
	svCGOReportHandler(parameter, asdu);
}

static void bindProxy(SVSubscriber subscriber, uintptr_t parameter){
	SVSubscriber_setListener(subscriber, svListenerProxy, (void *)parameter);
}
*/
import "C"
import (
	"sync"
	"sync/atomic"
	"unsafe"
)

type (
	SvSubscriberCallbackID uintptr

	SvSubscriberASDU struct {
		noCopy struct{}
		cAsdu  C.SVSubscriber_ASDU
	}

	SvReport struct {
		ReceiverASDU SvSubscriberASDU
	}

	SvReportHandler func(report *SvReport)

	SvReceiverConf struct {
		InterfaceID string
	}

	SvReceiver struct {
		cSvReceiver C.SVReceiver
		refs        map[SvSubscriberCallbackID]struct{}
	}
)

var (
	__svCallbackLocker struct {
		noCopy struct{}
		sync.RWMutex
		idOffset     atomic.Uintptr
		callbackRefs map[SvSubscriberCallbackID]struct {
			subscriber *SvSubscriber
			handler    SvReportHandler
		}
	}
)

func init() {
	__svCallbackLocker.Lock()
	defer __svCallbackLocker.Unlock()
	__svCallbackLocker.idOffset.Add(1000)
	__svCallbackLocker.callbackRefs = map[SvSubscriberCallbackID]struct {
		subscriber *SvSubscriber
		handler    SvReportHandler
	}{}
}

//export svCGOReportHandler
func svCGOReportHandler(parameter unsafe.Pointer, asdu C.SVSubscriber_ASDU) {
	objID := SvSubscriberCallbackID(parameter)
	__svCallbackLocker.RLock()
	defer __svCallbackLocker.RUnlock()

	if fetch, ok := __svCallbackLocker.callbackRefs[objID]; ok {
		fetch.handler(&SvReport{
			ReceiverASDU: SvSubscriberASDU{
				cAsdu: asdu,
			},
		})
	}
}

func NewSvReceiver(conf SvReceiverConf) *SvReceiver {
	ether := C.CString(conf.InterfaceID)
	defer C.free(unsafe.Pointer(ether))
	cSvReceiver := C.SVReceiver_create()
	C.SVReceiver_setInterfaceId(cSvReceiver, ether)

	return &SvReceiver{
		cSvReceiver: cSvReceiver,
		refs:        make(map[SvSubscriberCallbackID]struct{}),
	}
}

func (receiver *SvReceiver) AddSubscriber(subscriber *SvSubscriber) *SvReceiver {
	__svCallbackLocker.Lock()
	defer __svCallbackLocker.Unlock()

	__svCallbackLocker.callbackRefs[subscriber.callbackID] = struct {
		subscriber *SvSubscriber
		handler    SvReportHandler
	}{
		subscriber: subscriber,
		handler:    subscriber.conf.Handler,
	}
	C.bindProxy(subscriber.cSubscriber, C.uintptr_t(subscriber.callbackID))
	C.SVReceiver_addSubscriber(receiver.cSvReceiver, subscriber.cSubscriber)
	receiver.refs[subscriber.callbackID] = struct{}{}

	return receiver
}

func (receiver *SvReceiver) RemoveSubscriber(subscriber *SvSubscriber) *SvReceiver {
	__svCallbackLocker.Lock()
	defer __svCallbackLocker.Unlock()

	C.SVReceiver_removeSubscriber(receiver.cSvReceiver, subscriber.cSubscriber)
	delete(__svCallbackLocker.callbackRefs, subscriber.callbackID)
	delete(receiver.refs, subscriber.callbackID)

	return receiver
}

func (receiver *SvReceiver) IsRunning() bool {
	return bool(C.SVReceiver_isRunning(receiver.cSvReceiver))
}

func (receiver *SvReceiver) Tick() bool {
	return bool(C.SVReceiver_tick(receiver.cSvReceiver))
}

func (receiver *SvReceiver) DisableAddrCheck() *SvReceiver {
	C.SVReceiver_disableDestAddrCheck(receiver.cSvReceiver)

	return receiver
}

func (receiver *SvReceiver) EnableAddrCheck() *SvReceiver {
	C.SVReceiver_enableDestAddrCheck(receiver.cSvReceiver)

	return receiver
}

func (receiver *SvReceiver) Start() *SvReceiver {
	C.SVReceiver_start(receiver.cSvReceiver)

	return receiver
}

func (receiver *SvReceiver) Stop() *SvReceiver {
	C.SVReceiver_stop(receiver.cSvReceiver)

	return receiver
}

func (receiver *SvReceiver) Destroy() {
	__svCallbackLocker.Lock()
	defer __svCallbackLocker.Unlock()

	for id := range receiver.refs {
		delete(__svCallbackLocker.callbackRefs, id)
	}
	C.SVReceiver_destroy(receiver.cSvReceiver)

	receiver.refs = nil
	receiver.cSvReceiver = nil
}

func (receiver *SvSubscriberASDU) GetSmpCnt() uint16 {
	return uint16(C.SVSubscriber_ASDU_getSmpCnt(receiver.cAsdu))
}

func (receiver *SvSubscriberASDU) GetDatSet() string {
	return C.GoString(C.SVSubscriber_ASDU_getDatSet(receiver.cAsdu))
}

func (receiver *SvSubscriberASDU) GetSvID() string {
	return C.GoString(C.SVSubscriber_ASDU_getSvId(receiver.cAsdu))
}

func (receiver *SvSubscriberASDU) GetConfRev() uint32 {
	return uint32(C.SVSubscriber_ASDU_getConfRev(receiver.cAsdu))
}

func (receiver *SvSubscriberASDU) GetSmpMod() uint8 {
	return uint8(C.SVSubscriber_ASDU_getSmpMod(receiver.cAsdu))
}

func (receiver *SvSubscriberASDU) GetSmpRate() uint16 {
	return uint16(C.SVSubscriber_ASDU_getSmpRate(receiver.cAsdu))
}

func (receiver *SvSubscriberASDU) GetSmpSynch() uint8 {
	return uint8(C.SVSubscriber_ASDU_getSmpSynch(receiver.cAsdu))
}

func (receiver *SvSubscriberASDU) GetRefrTmAsMs() uint64 {
	return uint64(C.SVSubscriber_ASDU_getRefrTmAsMs(receiver.cAsdu))
}

func (receiver *SvSubscriberASDU) GetRefrTmAsNs() uint64 {
	return uint64(C.SVSubscriber_ASDU_getRefrTmAsNs(receiver.cAsdu))
}

func (receiver *SvSubscriberASDU) HasDatSet() bool {
	return bool(C.SVSubscriber_ASDU_hasDatSet(receiver.cAsdu))
}

func (receiver *SvSubscriberASDU) HasRefrTm() bool {
	return bool(C.SVSubscriber_ASDU_hasRefrTm(receiver.cAsdu))
}

func (receiver *SvSubscriberASDU) HasSmpMod() bool {
	return bool(C.SVSubscriber_ASDU_hasSmpMod(receiver.cAsdu))
}

func (receiver *SvSubscriberASDU) HasSmpRate() bool {
	return bool(C.SVSubscriber_ASDU_hasSmpRate(receiver.cAsdu))
}

func (receiver *SvSubscriberASDU) GetInt8(index int) int8 {
	return int8(C.SVSubscriber_ASDU_getINT8(receiver.cAsdu, C.int(index)))
}

func (receiver *SvSubscriberASDU) GetInt16(index int) int16 {
	return int16(C.SVSubscriber_ASDU_getINT16(receiver.cAsdu, C.int(index)))
}

func (receiver *SvSubscriberASDU) GetInt32(index int) int32 {
	return int32(C.SVSubscriber_ASDU_getINT32(receiver.cAsdu, C.int(index)))
}

func (receiver *SvSubscriberASDU) GetInt64(index int) int64 {
	return int64(C.SVSubscriber_ASDU_getINT64(receiver.cAsdu, C.int(index)))
}

func (receiver *SvSubscriberASDU) GetUint8(index int) uint8 {
	return uint8(C.SVSubscriber_ASDU_getINT8U(receiver.cAsdu, C.int(index)))
}

func (receiver *SvSubscriberASDU) GetUint16(index int) uint16 {
	return uint16(C.SVSubscriber_ASDU_getINT16U(receiver.cAsdu, C.int(index)))
}

func (receiver *SvSubscriberASDU) GetUint32(index int) uint32 {
	return uint32(C.SVSubscriber_ASDU_getINT32U(receiver.cAsdu, C.int(index)))
}

func (receiver *SvSubscriberASDU) GetUint64(index int) uint64 {
	return uint64(C.SVSubscriber_ASDU_getINT64U(receiver.cAsdu, C.int(index)))
}

func (receiver *SvSubscriberASDU) GetFloat32(index int) float32 {
	return float32(C.SVSubscriber_ASDU_getFLOAT32(receiver.cAsdu, C.int(index)))
}

func (receiver *SvSubscriberASDU) GetFloat64(index int) float64 {
	return float64(C.SVSubscriber_ASDU_getFLOAT64(receiver.cAsdu, C.int(index)))
}

func (receiver *SvSubscriberASDU) GetDataSize() int {
	return int(C.SVSubscriber_ASDU_getDataSize(receiver.cAsdu))
}

func (receiver *SvSubscriberASDU) GetTimestamp(index int) *Timestamp {
	cVal := C.SVSubscriber_ASDU_getTimestamp(receiver.cAsdu, C.int(index))
	return &Timestamp{
		cTimestamp: cVal,
	}
}

func (receiver *SvSubscriberASDU) GetQuality(index int) Quality {
	return Quality(C.SVSubscriber_ASDU_getQuality(receiver.cAsdu, C.int(index)))
}
