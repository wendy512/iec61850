//go:build linux && amd64

package iec61850

/*
#include "goose_receiver.h"
#include "goose_subscriber.h"
#include <stdint.h>

extern void cgoReportCallbackBridgeDispatcher(GooseSubscriber subscriber, void *parameter);

static void goose_report_proxy_handler(GooseSubscriber subscriber, void* parameter) {
	cgoReportCallbackBridgeDispatcher(subscriber, parameter);
}

static void simple_goose_subscriber_set_listener(GooseSubscriber subscriber, uintptr_t parameter) {
	GooseSubscriber_setListener(subscriber, goose_report_proxy_handler, (void *)parameter);
}
*/
import "C"
import (
	"sync"
	"sync/atomic"
	"unsafe"
)

type (
	GooseReceiver struct {
		noCopy        struct{}
		gooseReceiver *C.struct_sGooseReceiver
		refs          map[GooseCallbackHandlerID]struct{}
	}
)

var (
	gooseCallbackLocker struct {
		noCopy struct{}
		sync.RWMutex
		idOffset     atomic.Uintptr
		callbackRefs map[GooseCallbackHandlerID]struct {
			handler    GooseReportCallback
			subscriber *GooseSubscriber
		}
	}
)

func init() {
	gooseCallbackLocker.Lock()
	defer gooseCallbackLocker.Unlock()

	gooseCallbackLocker.idOffset.Add(1000)
	gooseCallbackLocker.callbackRefs = map[GooseCallbackHandlerID]struct {
		handler    GooseReportCallback
		subscriber *GooseSubscriber
	}{}
}

//export cgoReportCallbackBridgeDispatcher
func cgoReportCallbackBridgeDispatcher(_ *C.struct_sGooseSubscriber, parameter unsafe.Pointer) {
	refID := GooseCallbackHandlerID(parameter)
	gooseCallbackLocker.RLock()
	defer gooseCallbackLocker.RUnlock()

	if fetch, ok := gooseCallbackLocker.callbackRefs[refID]; ok {
		fetch.handler(&GooseReport{
			parameter:       parameter,
			GooseSubscriber: fetch.subscriber,
		})
	}
}

func NewGooseReceiver() *GooseReceiver {
	return &GooseReceiver{
		gooseReceiver: C.GooseReceiver_create(),
		refs:          make(map[GooseCallbackHandlerID]struct{}),
	}
}

func (receiver *GooseReceiver) AddSubscriber(subscriber *GooseSubscriber) *GooseReceiver {
	gooseCallbackLocker.Lock()
	defer gooseCallbackLocker.Unlock()

	gooseCallbackLocker.callbackRefs[subscriber.HandlerID] = struct {
		handler    GooseReportCallback
		subscriber *GooseSubscriber
	}{
		handler:    subscriber.Conf.ReportHandler,
		subscriber: subscriber,
	}
	receiver.refs[subscriber.HandlerID] = struct{}{}
	C.simple_goose_subscriber_set_listener(
		subscriber.subscriber,
		C.uintptr_t(subscriber.HandlerID),
	)
	C.GooseReceiver_addSubscriber(receiver.gooseReceiver, subscriber.subscriber)

	return receiver
}

func (receiver *GooseReceiver) RemoveSubscriber(subscriber *GooseSubscriber) *GooseReceiver {
	gooseCallbackLocker.Lock()
	defer gooseCallbackLocker.Unlock()

	C.GooseReceiver_removeSubscriber(receiver.gooseReceiver, subscriber.subscriber)
	delete(gooseCallbackLocker.callbackRefs, subscriber.HandlerID)
	delete(receiver.refs, subscriber.HandlerID)

	return receiver
}

func (receiver *GooseReceiver) SetInterfaceID(interfaceID string) *GooseReceiver {
	tmp := C.CString(interfaceID)
	defer C.free(unsafe.Pointer(tmp))
	C.GooseReceiver_setInterfaceId(receiver.gooseReceiver, tmp)

	return receiver
}

func (receiver *GooseReceiver) GetInterfaceID() string {
	return C.GoString(C.GooseReceiver_getInterfaceId(receiver.gooseReceiver))
}

func (receiver *GooseReceiver) Start() *GooseReceiver {
	C.GooseReceiver_start(receiver.gooseReceiver)

	return receiver
}

func (receiver *GooseReceiver) IsRunning() bool {
	return bool(C.GooseReceiver_isRunning(receiver.gooseReceiver))
}

func (receiver *GooseReceiver) Tick() bool {
	return bool(C.GooseReceiver_tick(receiver.gooseReceiver))
}

func (receiver *GooseReceiver) Stop() *GooseReceiver {
	C.GooseReceiver_stop(receiver.gooseReceiver)

	return receiver
}

func (receiver *GooseReceiver) Destroy() {
	gooseCallbackLocker.Lock()
	defer gooseCallbackLocker.Unlock()
	for id := range receiver.refs {
		delete(gooseCallbackLocker.callbackRefs, id)
	}
	C.GooseReceiver_destroy(receiver.gooseReceiver)
	receiver.refs = nil
	receiver.gooseReceiver = nil
}
