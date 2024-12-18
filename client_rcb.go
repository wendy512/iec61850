package iec61850

// #include <iec61850_client.h>
// #include <iec61850_common.h>
import "C"

import (
	"unsafe"
)

type TrgOps struct {
	DataChange            bool // 值变化
	QualityChange         bool // 品质变化
	DataUpdate            bool // 数据修改
	TriggeredPeriodically bool // 周期触发
	Gi                    bool // GI(一般审问)请求触发
	Transient             bool // 瞬变
}
type OptFlds struct {
	SequenceNumber     bool // 顺序号
	TimeOfEntry        bool // 报告时标
	ReasonForInclusion bool // 原因码
	DataSetName        bool // 数据集
	DataReference      bool // 数据引用
	BufferOverflow     bool // 缓存溢出标识
	EntryID            bool // 报告标识符
	ConfigRevision     bool // 配置版本号
}

type ClientReportControlBlock struct {
	Ena     bool    // 使能
	IntgPd  int     // 周期上送时间
	TrgOps  TrgOps  // 触发条件
	OptFlds OptFlds // 报告选项
}

func (c *Client) GetRCBValues(objectReference string) (*ClientReportControlBlock, error) {
	var clientError C.IedClientError
	cObjectRef := C.CString(objectReference)
	defer C.free(unsafe.Pointer(cObjectRef))
	rcb := C.IedConnection_getRCBValues(c.conn, &clientError, cObjectRef, nil)
	if rcb == nil {
		return nil, GetIedClientError(clientError)
	}
	return &ClientReportControlBlock{
		Ena:     c.getRCBEnable(rcb),
		IntgPd:  int(c.getRCBIntgPd(rcb)),
		TrgOps:  c.getTrgOps(rcb),
		OptFlds: c.getOptFlds(rcb),
	}, nil
}

func (c *Client) getRCBEnable(rcb C.ClientReportControlBlock) bool {
	enable := C.ClientReportControlBlock_getRptEna(rcb)
	return bool(enable)
}

func (c *Client) getRCBIntgPd(rcb C.ClientReportControlBlock) uint32 {
	intgPd := C.ClientReportControlBlock_getIntgPd(rcb)
	return uint32(intgPd)
}

func (c *Client) getOptFlds(rcb C.ClientReportControlBlock) OptFlds {
	optFlds := C.ClientReportControlBlock_getOptFlds(rcb)
	g := int(optFlds)
	return OptFlds{
		SequenceNumber:     IsBitSet(g, 0),
		TimeOfEntry:        IsBitSet(g, 1),
		ReasonForInclusion: IsBitSet(g, 2),
		DataSetName:        IsBitSet(g, 3),
		DataReference:      IsBitSet(g, 4),
		BufferOverflow:     IsBitSet(g, 5),
		EntryID:            IsBitSet(g, 6),
		ConfigRevision:     IsBitSet(g, 7),
	}
}

func (c *Client) getTrgOps(rcb C.ClientReportControlBlock) TrgOps {
	trgOps := C.ClientReportControlBlock_getTrgOps(rcb)
	g := int(trgOps)
	return TrgOps{
		DataChange:            IsBitSet(g, 0),
		QualityChange:         IsBitSet(g, 1),
		DataUpdate:            IsBitSet(g, 2),
		TriggeredPeriodically: IsBitSet(g, 3),
		Gi:                    IsBitSet(g, 4),
		Transient:             IsBitSet(g, 5),
	}
}

func (c *Client) SetRCBValues(objectReference string, settings ClientReportControlBlock) error {
	var clientError C.IedClientError
	cObjectRef := C.CString(objectReference)
	defer C.free(unsafe.Pointer(cObjectRef))
	rcb := C.ClientReportControlBlock_create(cObjectRef)
	defer C.ClientReportControlBlock_destroy(rcb)
	var trgOps, optFlds C.int
	// trgOps
	if settings.TrgOps.DataChange {
		trgOps = trgOps | C.TRG_OPT_DATA_CHANGED
	}
	if settings.TrgOps.QualityChange {
		trgOps = trgOps | C.TRG_OPT_QUALITY_CHANGED
	}
	if settings.TrgOps.DataUpdate {
		trgOps = trgOps | C.TRG_OPT_DATA_UPDATE
	}
	if settings.TrgOps.TriggeredPeriodically {
		trgOps = trgOps | C.TRG_OPT_INTEGRITY
	}
	if settings.TrgOps.Gi {
		trgOps = trgOps | C.TRG_OPT_GI
	}
	if settings.TrgOps.Transient {
		trgOps = trgOps | C.TRG_OPT_TRANSIENT
	}
	// optFlds
	if settings.OptFlds.SequenceNumber {
		optFlds = optFlds | C.RPT_OPT_SEQ_NUM
	}
	if settings.OptFlds.TimeOfEntry {
		optFlds = optFlds | C.RPT_OPT_TIME_STAMP
	}
	if settings.OptFlds.ReasonForInclusion {
		optFlds = optFlds | C.RPT_OPT_REASON_FOR_INCLUSION
	}
	if settings.OptFlds.DataSetName {
		optFlds = optFlds | C.RPT_OPT_DATA_SET
	}
	if settings.OptFlds.DataReference {
		optFlds = optFlds | C.RPT_OPT_DATA_REFERENCE
	}
	if settings.OptFlds.BufferOverflow {
		optFlds = optFlds | C.RPT_OPT_BUFFER_OVERFLOW
	}
	if settings.OptFlds.EntryID {
		optFlds = optFlds | C.RPT_OPT_ENTRY_ID
	}
	if settings.OptFlds.ConfigRevision {
		optFlds = optFlds | C.RPT_OPT_CONF_REV
	}

	C.ClientReportControlBlock_setTrgOps(rcb, trgOps)                      // 触发条件
	C.ClientReportControlBlock_setRptEna(rcb, C.bool(settings.Ena))        // 报告使能
	C.ClientReportControlBlock_setIntgPd(rcb, C.uint32_t(settings.IntgPd)) // 周期上送时间
	C.ClientReportControlBlock_setOptFlds(rcb, optFlds)
	C.IedConnection_setRCBValues(c.conn, &clientError, rcb, C.RCB_ELEMENT_RPT_ENA|C.RCB_ELEMENT_TRG_OPS|C.RCB_ELEMENT_INTG_PD, true)
	if err := GetIedClientError(clientError); err != nil {
		return err
	}
	return nil
}

func IsBitSet(val int, pos int) bool {
	return (val & (1 << pos)) != 0
}
