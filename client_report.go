package iec61850

/*
#include <iec61850_client.h>

extern void reportCallbackFunctionBridge(void* parameter, ClientReport report);
*/
import "C"
import (
	"unsafe"
)

var reportCallbacks = make(map[int32]*reportCallbackHandler)

type ReasonForInclusion int

const (
	IEC61850_REASON_NOT_INCLUDED ReasonForInclusion = iota
	IEC61850_REASON_DATA_CHANGE
	IEC61850_REASON_QUALITY_CHANGE
	IEC61850_REASON_DATA_UPDATE
	IEC61850_REASON_INTEGRITY
	IEC61850_REASON_GI
	IEC61850_REASON_UNKNOWN
)

type reportCallbackHandler struct {
	handler ReportCallbackFunction
}

type ClientReport struct {
	Report C.ClientReport
}

type ReportCallbackFunction func(clientReport ClientReport)

//export reportCallbackFunctionBridge
func reportCallbackFunctionBridge(parameter unsafe.Pointer, report C.ClientReport) {
	callbackId := int32(uintptr(parameter))
	if call, ok := reportCallbacks[callbackId]; ok {
		call.handler(ClientReport{
			Report: report,
		})
	}
}

func (reason ReasonForInclusion) GetValueAsString() string {
	return C.GoString(C.ReasonForInclusion_getValueAsString(C.ReasonForInclusion(reason)))
}

func (clientReport *ClientReport) GetElement(elementIndex int) (MmsValue, error) {
	dataSetValues := C.ClientReport_getDataSetValues(clientReport.Report)
	cMmsValue := C.MmsValue_getElement(dataSetValues, C.int(elementIndex))
	mmsType := MmsType(C.MmsValue_getType(cMmsValue))

	MmsVal, err := toGoValue(cMmsValue, mmsType)
	if err != nil {
		return MmsValue{}, err
	}

	return MmsValue{
		Value: MmsVal,
		Type:  mmsType,
	}, nil
}

func (clientReport *ClientReport) GetRcbReference() string {
	return C.GoString(C.ClientReport_getRcbReference(clientReport.Report))
}

func (clientReport *ClientReport) GetRptId() string {
	return C.GoString(C.ClientReport_getRptId(clientReport.Report))
}

func (clientReport *ClientReport) GetReasonForInclusion(elementIndex int) ReasonForInclusion {
	reason := C.ClientReport_getReasonForInclusion(clientReport.Report, C.int(elementIndex))
	return ReasonForInclusion(reason)
}

func (clientReport *ClientReport) HasTimestamp() bool {
	return bool(C.ClientReport_hasTimestamp(clientReport.Report))
}

func (clientReport *ClientReport) GetTimestamp() int64 {
	return int64(C.ClientReport_getTimestamp(clientReport.Report))
}

func (clientReport *ClientReport) HasSeqNum() bool {
	return bool(C.ClientReport_hasSeqNum(clientReport.Report))
}

func (clientReport *ClientReport) GetSeqNum() int16 {
	return int16(C.ClientReport_getSeqNum(clientReport.Report))
}

func (clientReport *ClientReport) HasDataSetName() bool {
	return bool(C.ClientReport_hasDataSetName(clientReport.Report))
}

func (clientReport *ClientReport) HasReasonForInclusion() bool {
	return bool(C.ClientReport_hasReasonForInclusion(clientReport.Report))
}

func (clientReport *ClientReport) HasConfRev() bool {
	return bool(C.ClientReport_hasConfRev(clientReport.Report))
}

func (clientReport *ClientReport) GetConfRev() int32 {
	return int32(C.ClientReport_getConfRev(clientReport.Report))
}

func (clientReport *ClientReport) HasBufOvfl() bool {
	return bool(C.ClientReport_hasBufOvfl(clientReport.Report))
}

func (clientReport *ClientReport) GetBufOvfl() bool {
	return bool(C.ClientReport_getBufOvfl(clientReport.Report))
}

func (clientReport *ClientReport) HasDataReference() bool {
	return bool(C.ClientReport_hasDataReference(clientReport.Report))
}

func (clientReport *ClientReport) GetDataReference(elementIndex int) string {
	return C.GoString(C.ClientReport_getDataReference(clientReport.Report, C.int(elementIndex)))
}

func (clientReport *ClientReport) GetDataSetName() string {
	return C.GoString(C.ClientReport_getDataSetName(clientReport.Report))
}

func (clientReport *ClientReport) GetDataSetValues() (MmsValue, error) {
	cMmsValue := C.ClientReport_getDataSetValues(clientReport.Report)
	mmsType := MmsType(C.MmsValue_getType(cMmsValue))

	mmsValue, err := toGoValue(cMmsValue, mmsType)
	if err != nil {
		return MmsValue{}, err
	}

	return MmsValue{
		Value: mmsValue,
		Type:  mmsType,
	}, nil
}

func (clientReport *ClientReport) HasSubSeqNum() bool {
	return bool(C.ClientReport_hasSubSeqNum(clientReport.Report))
}

func (clientReport *ClientReport) GetSubSeqNum() int16 {
	return int16(C.ClientReport_getSubSeqNum(clientReport.Report))
}

func (clientReport *ClientReport) GetMoreSeqmentsFollow() bool {
	return bool(C.ClientReport_getMoreSeqmentsFollow(clientReport.Report))
}

func (c *Client) InstallReportHandler(objectReference string, function ReportCallbackFunction) error {
	var clientError C.IedClientError

	cObjectRef := C.CString(objectReference)
	defer C.free(unsafe.Pointer(cObjectRef))

	rcb := C.IedConnection_getRCBValues(c.conn, &clientError, cObjectRef, nil)
	if err := GetIedClientError(clientError); err != nil {
		return err
	}
	defer C.ClientReportControlBlock_destroy(rcb)

	callbackId := callbackIdGen.Add(1)
	cPtr := intToPointerBug58625(callbackId)
	reportCallbacks[callbackId] = &reportCallbackHandler{
		handler: function,
	}

	C.IedConnection_installReportHandler(c.conn, cObjectRef, C.ClientReportControlBlock_getRptId(rcb), (*[0]byte)(C.reportCallbackFunctionBridge), cPtr)

	return nil
}

func (c *Client) UninstallReportHandler(objectReference string) {
	cObjectRef := C.CString(objectReference)
	defer C.free(unsafe.Pointer(cObjectRef))
	C.IedConnection_uninstallReportHandler(c.conn, cObjectRef)
}

func (c *Client) TriggerGIReport(objectReference string) error {
	var clientError C.IedClientError

	cObjectRef := C.CString(objectReference)
	defer C.free(unsafe.Pointer(cObjectRef))

	rcb := C.IedConnection_getRCBValues(c.conn, &clientError, cObjectRef, nil)
	if err := GetIedClientError(clientError); err != nil {
		return err
	}
	defer C.ClientReportControlBlock_destroy(rcb)

	C.IedConnection_triggerGIReport(c.conn, &clientError, cObjectRef)
	if err := GetIedClientError(clientError); err != nil {
		return err
	}

	return nil
}
