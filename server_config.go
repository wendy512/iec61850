package iec61850

// #include <iec61850_server.h>
import "C"
import "unsafe"

// ServerConfig Configuration object to configure IEC 61850 stack features
type ServerConfig struct {
	Edition                        uint8  // IEC 61850 edition (0 = edition 1, 1 = edition 2, 2 = edition 2.1, ...)
	ReportBufferSize               int    // size of the report buffer associated with a buffered report control block
	ReportBufferSizeForURCBs       int    // size of the report buffer associated with an unbuffered report control block
	MaxConnections                 int    // maximum number of MMS (TCP) connections
	SyncIntegrityReportTimes       bool   // integrity report start times will by synchronized with straight numbers
	EnableFileService              bool   // when true (default) enable MMS file service
	FileServiceBasePath            string // Base path (directory where the file service serves files
	EnableDynamicDataSetService    bool   // when true (default) enable dynamic data set services for MMS
	MaxAssociationSpecificDataSets int    // the maximum number of allowed association specific data sets
	MaxDomainSpecificDataSets      int    // the maximum number of allowed domain specific data sets
	MaxDataSetEntries              int    // maximum number of data set entries of dynamic data sets
	EnableLogService               bool   // when true (default) enable log service
	EnableEditSG                   bool   // enable EditSG service
	EnableResvTmsForSGCB           bool   // enable visibility of SGCB.ResvTms
	EnableResvTmsForBRCB           bool   // BRCB has resvTms attribute - only edition 2
	EnableOwnerForRCB              bool   // RCB has owner attribute
	UseIntegratedGoosePublisher    bool   // when true (default) the integrated GOOSE publisher is used
	reportSettings                 ReportSetting
}

type ReportSetting struct {
	setting uint8
	isDyn   bool
}

// NewServerConfig creates a new ServerConfig object with default values
func NewServerConfig() ServerConfig {
	return ServerConfig{
		Edition:                        1,
		ReportBufferSize:               65536,
		ReportBufferSizeForURCBs:       65536,
		MaxConnections:                 5,
		SyncIntegrityReportTimes:       false,
		EnableFileService:              true,
		EnableDynamicDataSetService:    true,
		FileServiceBasePath:            "./vmd-filestore/",
		MaxAssociationSpecificDataSets: 10,
		MaxDomainSpecificDataSets:      10,
		MaxDataSetEntries:              100,
		EnableLogService:               true,
		EnableEditSG:                   true,
		EnableResvTmsForSGCB:           true,
		EnableResvTmsForBRCB:           true,
		EnableOwnerForRCB:              false,
		UseIntegratedGoosePublisher:    true,
	}
}

// SetReportSetting Make a configurable report setting writeable or read-only
//
// Parameters:
//
//	setting: one of IEC61850_REPORTSETTINGS_RPT_ID, _BUF_TIME, _DATSET, _TRG_OPS, _OPT_FIELDS, _INTG_PD
//	isDyn: true, when setting is writable ("Dyn") or false, when read-only
func (that ServerConfig) SetReportSetting(setting uint8, isDyn bool) {
	that.reportSettings = ReportSetting{
		setting: setting,
		isDyn:   isDyn,
	}
}

func (that ServerConfig) createIedServerConfig(serverConfig ServerConfig) C.IedServerConfig {
	config := C.IedServerConfig_create()

	cFileServiceBasePath := C.CString(serverConfig.FileServiceBasePath)
	defer C.free(unsafe.Pointer(cFileServiceBasePath))

	C.IedServerConfig_setEdition(config, C.uint8_t(serverConfig.Edition))
	C.IedServerConfig_setReportBufferSize(config, C.int(serverConfig.ReportBufferSize))
	C.IedServerConfig_setReportBufferSizeForURCBs(config, C.int(serverConfig.ReportBufferSizeForURCBs))
	C.IedServerConfig_setMaxMmsConnections(config, C.int(serverConfig.MaxConnections))
	C.IedServerConfig_setSyncIntegrityReportTimes(config, C.bool(serverConfig.SyncIntegrityReportTimes))
	C.IedServerConfig_enableFileService(config, C.bool(serverConfig.EnableFileService))
	C.IedServerConfig_setFileServiceBasePath(config, cFileServiceBasePath)
	C.IedServerConfig_enableDynamicDataSetService(config, C.bool(serverConfig.EnableDynamicDataSetService))
	C.IedServerConfig_setMaxAssociationSpecificDataSets(config, C.int(serverConfig.MaxAssociationSpecificDataSets))
	C.IedServerConfig_setMaxDomainSpecificDataSets(config, C.int(serverConfig.MaxDomainSpecificDataSets))
	C.IedServerConfig_setMaxDataSetEntries(config, C.int(serverConfig.MaxDataSetEntries))
	C.IedServerConfig_enableLogService(config, C.bool(serverConfig.EnableLogService))
	C.IedServerConfig_enableEditSG(config, C.bool(serverConfig.EnableEditSG))
	C.IedServerConfig_enableResvTmsForSGCB(config, C.bool(serverConfig.EnableResvTmsForSGCB))
	C.IedServerConfig_enableResvTmsForBRCB(config, C.bool(serverConfig.EnableResvTmsForBRCB))
	C.IedServerConfig_enableOwnerForRCB(config, C.bool(serverConfig.EnableOwnerForRCB))
	C.IedServerConfig_useIntegratedGoosePublisher(config, C.bool(serverConfig.UseIntegratedGoosePublisher))
	C.IedServerConfig_setReportSetting(config, C.uint8_t(serverConfig.reportSettings.setting), C.bool(serverConfig.reportSettings.isDyn))
	return config
}
