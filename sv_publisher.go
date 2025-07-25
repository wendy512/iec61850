package iec61850

/*
#include "sv_publisher.h"
*/
import "C"
import (
	"sync"
	"unsafe"
)

type (
	SvPublisherConf struct {
		EtherName    string
		AppID        uint16
		DstAddr      [6]uint8
		VlanID       uint16
		VlanPriority uint8
	}

	SVPublisher struct {
		noCopy struct{}
		sync.RWMutex
		cSvPublisher C.SVPublisher
		Conf         SvPublisherConf
		postCleanup  []unsafe.Pointer
	}
)

func NewSVPublisher(conf SvPublisherConf) *SVPublisher {
	parameters := C.struct_sCommParameters{}
	parameters.appId = C.uint16_t(conf.AppID)
	parameters.vlanId = C.uint16_t(conf.VlanID)
	parameters.vlanPriority = C.uint8_t(conf.VlanPriority)
	for i := 0; i < len(conf.DstAddr); i++ {
		parameters.dstAddress[i] = C.uint8_t(conf.DstAddr[i])
	}
	ether := C.CString(conf.EtherName)
	defer C.free(unsafe.Pointer(ether))
	cSvPublisher := C.SVPublisher_create(&parameters, ether)

	return &SVPublisher{
		cSvPublisher: cSvPublisher,
		Conf:         conf,
	}
}

func (receiver *SVPublisher) AddSvASDU(svID string, dataset string, confRev uint32) *SvPublisherASDU {
	cSvID := C.CString(svID)
	cDataset := C.CString(dataset)

	receiver.Lock()
	defer receiver.Unlock()
	receiver.postCleanup = append(receiver.postCleanup, unsafe.Pointer(cSvID), unsafe.Pointer(cDataset))

	v := C.SVPublisher_addASDU(receiver.cSvPublisher, cSvID, cDataset, C.uint32_t(confRev))
	return &SvPublisherASDU{
		cPublisherASDU: v,
	}
}

func (receiver *SVPublisher) SetupComplete() *SVPublisher {
	C.SVPublisher_setupComplete(receiver.cSvPublisher)
	return receiver
}

func (receiver *SVPublisher) Publish() {
	receiver.RLock()
	defer receiver.RUnlock()

	C.SVPublisher_publish(receiver.cSvPublisher)
}

func (receiver *SVPublisher) Destroy() {
	C.SVPublisher_destroy(receiver.cSvPublisher)
	for _, ptr := range receiver.postCleanup {
		C.free(ptr)
	}
	receiver.postCleanup = nil
	receiver.cSvPublisher = nil
}

type (
	SvPublisherASDU struct {
		noCopy         struct{}
		cPublisherASDU C.SVPublisher_ASDU
	}
)

func (receiver *SvPublisherASDU) ResetBuffer() {
	C.SVPublisher_ASDU_resetBuffer(receiver.cPublisherASDU)
}

func (receiver *SvPublisherASDU) AddInt8() int {
	return int(C.SVPublisher_ASDU_addINT8(receiver.cPublisherASDU))
}

func (receiver *SvPublisherASDU) SetInt8(index int, value int8) {
	C.SVPublisher_ASDU_setINT8(receiver.cPublisherASDU, C.int(index), C.int8_t(value))
}

func (receiver *SvPublisherASDU) AddInt32() int {
	return int(C.SVPublisher_ASDU_addINT32(receiver.cPublisherASDU))
}

func (receiver *SvPublisherASDU) SetInt32(index int, value int32) {
	C.SVPublisher_ASDU_setINT32(receiver.cPublisherASDU, C.int(index), C.int32_t(value))
}

func (receiver *SvPublisherASDU) AddInt64() int {
	return int(C.SVPublisher_ASDU_addINT64(receiver.cPublisherASDU))
}

func (receiver *SvPublisherASDU) SetInt64(index int, value int64) {
	C.SVPublisher_ASDU_setINT64(receiver.cPublisherASDU, C.int(index), C.int64_t(value))
}

func (receiver *SvPublisherASDU) AddFloat64() int {
	return int(C.SVPublisher_ASDU_addFLOAT64(receiver.cPublisherASDU))
}

func (receiver *SvPublisherASDU) SetFloat64(index int, value float64) {
	C.SVPublisher_ASDU_setFLOAT64(receiver.cPublisherASDU, C.int(index), C.double(value))
}

func (receiver *SvPublisherASDU) SetFloat(index int, value float32) {
	C.SVPublisher_ASDU_setFLOAT(receiver.cPublisherASDU, C.int(index), C.float(value))
}

func (receiver *SvPublisherASDU) AddFloat() int {
	return int(C.SVPublisher_ASDU_addFLOAT(receiver.cPublisherASDU))
}

func (receiver *SvPublisherASDU) AddTimestamp() int {
	return int(C.SVPublisher_ASDU_addTimestamp(receiver.cPublisherASDU))
}

func (receiver *SvPublisherASDU) SetTimestamp(index int, ts *Timestamp) {
	C.SVPublisher_ASDU_setTimestamp(receiver.cPublisherASDU, C.int(index), ts.cTimestamp)
}

func (receiver *SvPublisherASDU) AddQuality() int {
	return int(C.SVPublisher_ASDU_addQuality(receiver.cPublisherASDU))
}

func (receiver *SvPublisherASDU) SetQuality(index int, quality Quality) {
	C.SVPublisher_ASDU_setQuality(receiver.cPublisherASDU, C.int(index), C.Quality(quality))
}

func (receiver *SvPublisherASDU) SetSmpCnt(value uint16) {
	C.SVPublisher_ASDU_setSmpCnt(receiver.cPublisherASDU, C.uint16_t(value))
}

func (receiver *SvPublisherASDU) GetSmpCnt() uint16 {
	return uint16(C.SVPublisher_ASDU_getSmpCnt(receiver.cPublisherASDU))
}

func (receiver *SvPublisherASDU) IncreaseSmpCnt() {
	C.SVPublisher_ASDU_increaseSmpCnt(receiver.cPublisherASDU)
}

func (receiver *SvPublisherASDU) SetSmpCntWrap(limit uint16) {
	C.SVPublisher_ASDU_setSmpCntWrap(receiver.cPublisherASDU, C.uint16_t(limit))
}

func (receiver *SvPublisherASDU) EnableRefrTm() {
	C.SVPublisher_ASDU_enableRefrTm(receiver.cPublisherASDU)
}

func (receiver *SvPublisherASDU) SetRefrTmNs(nano int64) {
	C.SVPublisher_ASDU_setRefrTmNs(receiver.cPublisherASDU, C.nsSinceEpoch(nano))
}

func (receiver *SvPublisherASDU) SetRefrTmMs(ms int64) {
	C.SVPublisher_ASDU_setRefrTm(receiver.cPublisherASDU, C.msSinceEpoch(ms))
}

func (receiver *SvPublisherASDU) SetRefrTmByTimestamp(ts Timestamp) {
	C.SVPublisher_ASDU_setRefrTmByTimestamp(receiver.cPublisherASDU, &ts.cTimestamp)
}

func (receiver *SvPublisherASDU) SetSmpMod(mode uint8) {
	C.SVPublisher_ASDU_setSmpMod(receiver.cPublisherASDU, C.uint8_t(mode))
}

func (receiver *SvPublisherASDU) SetSmpRate(rate uint16) {
	C.SVPublisher_ASDU_setSmpRate(receiver.cPublisherASDU, C.uint16_t(rate))
}

func (receiver *SvPublisherASDU) SetSmpSynch(sync uint16) {
	C.SVPublisher_ASDU_setSmpSynch(receiver.cPublisherASDU, C.uint16_t(sync))
}
