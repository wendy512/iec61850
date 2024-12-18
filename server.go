package iec61850

// #include <iec61850_server.h>
import "C"
import (
	"unsafe"
)

type IedServer struct {
	server              C.IedServer
	serverConfig        ServerConfig
	tlsConfig           C.TLSConfiguration
	clientAuthenticator ClientAuthenticator
}

func NewServerWithTlsSupport(serverConfig ServerConfig, tlsConfig *TLSConfig, iedModel *IedModel) (*IedServer, error) {
	cTlsConfig, err := tlsConfig.createCTlsConfig()
	if err != nil {
		return nil, err
	}

	config := serverConfig.createIedServerConfig(serverConfig)
	defer C.IedServerConfig_destroy(config)
	return &IedServer{
		server:       C.IedServer_createWithConfig(iedModel._iedModel, cTlsConfig, config),
		serverConfig: serverConfig,
		tlsConfig:    cTlsConfig,
	}, nil
}

func NewServerWithConfig(serverConfig ServerConfig, iedModel *IedModel) *IedServer {
	config := serverConfig.createIedServerConfig(serverConfig)
	defer C.IedServerConfig_destroy(config)
	return &IedServer{
		server:       C.IedServer_createWithConfig(iedModel._iedModel, nil, config),
		serverConfig: serverConfig,
	}
}

// NewServer creates a new instance of the IedServer using the provided _iedModel.
func NewServer(iedModel *IedModel) *IedServer {
	return &IedServer{
		server: C.IedServer_create(iedModel._iedModel),
	}
}

// Start initiates the IedServer on the provided port.
func (is *IedServer) Start(port int) {
	C.IedServer_start(is.server, C.int(port))
	// If there's another way to detect the error, handle it here.
}

// IsRunning checks if the IedServer is currently running.
func (is *IedServer) IsRunning() bool {
	return bool(C.IedServer_isRunning(is.server))
}

// Stop terminates the IedServer.
func (is *IedServer) Stop() {
	C.IedServer_stop(is.server)
}

// Destroy frees all resources associated with the IedServer.
func (is *IedServer) Destroy() {
	C.IedServer_destroy(is.server)
}

// LockDataModel locks the data _iedModel of the IedServer.
func (is *IedServer) LockDataModel() {
	C.IedServer_lockDataModel(is.server)
}

// UnlockDataModel unlocks the data _iedModel of the IedServer.
func (is *IedServer) UnlockDataModel() {
	C.IedServer_unlockDataModel(is.server)
}

// UpdateUTCTimeAttributeValue updates a DataAttribute with a UTC time value.
func (is *IedServer) UpdateUTCTimeAttributeValue(node *ModelNode, value int64) {
	if node == nil || node._modelNode == nil {
		return
	}
	C.IedServer_updateUTCTimeAttributeValue(is.server, (*C.DataAttribute)(node._modelNode), C.uint64_t(value))
}

// UpdateFloatAttributeValue updates a DataAttribute with a float value.
func (is *IedServer) UpdateFloatAttributeValue(node *ModelNode, value float32) {
	if node == nil || node._modelNode == nil {
		return
	}
	C.IedServer_updateFloatAttributeValue(is.server, (*C.DataAttribute)(node._modelNode), C.float(value))
}

// UpdateInt32AttributeValue updates a DataAttribute with an Int32 value.
func (is *IedServer) UpdateInt32AttributeValue(node *ModelNode, value int32) {
	if node == nil || node._modelNode == nil {
		return
	}
	C.IedServer_updateInt32AttributeValue(is.server, (*C.DataAttribute)(node._modelNode), C.int32_t(value))
}

// UpdateVisibleStringAttributeValue updates a DataAttribute with a visible string value.
func (is *IedServer) UpdateVisibleStringAttributeValue(attr *DataAttribute, value string) {
	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))
	C.IedServer_updateVisibleStringAttributeValue(is.server, attr.attribute, cValue)
}
