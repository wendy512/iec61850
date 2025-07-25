//go:build linux && amd64

package iec61850

/*
#include "sv_subscriber.h"
*/
import "C"
import (
	"unsafe"
)

type (
	SvSubscriberConf struct {
		EthAddr [6]uint8
		AppID   uint16
		Handler SvReportHandler
	}

	SvSubscriber struct {
		noCopy      struct{}
		cSubscriber C.SVSubscriber
		callbackID  SvSubscriberCallbackID
		conf        SvSubscriberConf
	}
)

func NewSvSubscriber(conf SvSubscriberConf) *SvSubscriber {
	macAddr := (*C.uint8_t)(unsafe.Pointer(&conf.EthAddr[0]))
	subscriber := C.SVSubscriber_create(macAddr, C.uint16_t(conf.AppID))
	objID := SvSubscriberCallbackID(__svCallbackLocker.idOffset.Add(1))

	sub := &SvSubscriber{
		cSubscriber: subscriber,
		callbackID:  objID,
		conf:        conf,
	}

	return sub
}
