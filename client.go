package iec61850

// #include <iec61850_client.h>
import "C"
import (
	"sync/atomic"
	"unsafe"
)

type Client struct {
	conn      C.IedConnection
	connected *atomic.Bool
}

// Settings 连接配置
type Settings struct {
	Host           string
	Port           int
	ConnectTimeout uint // 连接超时配置，单位：毫秒
	RequestTimeout uint // 请求超时配置，单位：毫秒
}

func NewSettings() *Settings {
	return &Settings{
		Host:           "localhost",
		Port:           102,
		ConnectTimeout: 10000,
		RequestTimeout: 10000,
	}
}

// NewClient 创建客户端实例
func NewClient(settings *Settings) (*Client, error) {
	conn, clientErr := connect(settings)
	if err := GetIedClientError(clientErr); err != nil {
		return nil, err
	}

	connected := &atomic.Bool{}
	connected.Store(true)
	connection := &Client{
		conn:      conn,
		connected: connected,
	}
	return connection, nil
}

// Write 写单个属性值，不支持Structure
func (c *Client) Write(objectRef string, fc FC, value interface{}) error {
	mmsType, err := c.GetVariableSpecType(objectRef, fc)
	if err != nil {
		return err
	}
	var (
		mmsValue    *C.MmsValue
		clientError C.IedClientError
	)

	mmsValue, err = toMmsValue(mmsType, value)
	if err != nil {
		return err
	}

	cObjectRef := C.CString(objectRef)
	defer C.free(unsafe.Pointer(cObjectRef))
	defer C.MmsValue_delete(mmsValue)
	C.IedConnection_writeObject(c.conn, &clientError, cObjectRef, C.FunctionalConstraint(fc), mmsValue)
	return GetIedClientError(clientError)
}

// ReadBool 读取bool类型值
func (c *Client) ReadBool(objectRef string, fc FC) (bool, error) {
	cObjectRef := C.CString(objectRef)
	defer C.free(unsafe.Pointer(cObjectRef))

	var clientError C.IedClientError
	value := C.IedConnection_readBooleanValue(c.conn, &clientError, cObjectRef, C.FunctionalConstraint(fc))
	if err := GetIedClientError(clientError); err != nil {
		return false, err
	}
	return bool(value), nil
}

// ReadInt32 读取int32类型值
func (c *Client) ReadInt32(objectRef string, fc FC) (int32, error) {
	cObjectRef := C.CString(objectRef)
	defer C.free(unsafe.Pointer(cObjectRef))

	var clientError C.IedClientError
	value := C.IedConnection_readInt32Value(c.conn, &clientError, cObjectRef, C.FunctionalConstraint(fc))
	if err := GetIedClientError(clientError); err != nil {
		return 0, err
	}
	return int32(value), nil
}

// ReadInt64 读取int64类型值
func (c *Client) ReadInt64(objectRef string, fc FC) (int64, error) {
	cObjectRef := C.CString(objectRef)
	defer C.free(unsafe.Pointer(cObjectRef))

	var clientError C.IedClientError
	value := C.IedConnection_readInt64Value(c.conn, &clientError, cObjectRef, C.FunctionalConstraint(fc))
	if err := GetIedClientError(clientError); err != nil {
		return 0, err
	}
	return int64(value), nil
}

// ReadUint32 读取uint32类型值
func (c *Client) ReadUint32(objectRef string, fc FC) (uint32, error) {
	cObjectRef := C.CString(objectRef)
	defer C.free(unsafe.Pointer(cObjectRef))

	var clientError C.IedClientError
	value := C.IedConnection_readUnsigned32Value(c.conn, &clientError, cObjectRef, C.FunctionalConstraint(fc))
	if err := GetIedClientError(clientError); err != nil {
		return 0, err
	}
	return uint32(value), nil
}

// ReadFloat 读取float类型值
func (c *Client) ReadFloat(objectRef string, fc FC) (float32, error) {
	cObjectRef := C.CString(objectRef)
	defer C.free(unsafe.Pointer(cObjectRef))

	var clientError C.IedClientError
	value := C.IedConnection_readFloatValue(c.conn, &clientError, cObjectRef, C.FunctionalConstraint(fc))
	if err := GetIedClientError(clientError); err != nil {
		return 0, err
	}
	//源码返回值是C的float，4byte，所以应返回float32，否则会出现其他问题
	return float32(value), nil
}

// ReadString 读取string类型值
func (c *Client) ReadString(objectRef string, fc FC) (string, error) {
	cObjectRef := C.CString(objectRef)
	defer C.free(unsafe.Pointer(cObjectRef))

	var clientError C.IedClientError
	value := C.IedConnection_readStringValue(c.conn, &clientError, cObjectRef, C.FunctionalConstraint(fc))
	if err := GetIedClientError(clientError); err != nil {
		return "", err
	}
	return C.GoString(value), nil
}

// Read 读取属性数据
func (c *Client) Read(objectRef string, fc FC) (interface{}, error) {
	var clientError C.IedClientError
	cObjectRef := C.CString(objectRef)
	defer C.free(unsafe.Pointer(cObjectRef))

	mmsValue := C.IedConnection_readObject(c.conn, &clientError, cObjectRef, C.FunctionalConstraint(fc))
	if err := GetIedClientError(clientError); err != nil {
		return nil, err
	}

	defer C.MmsValue_delete(mmsValue)
	mmsType := MmsType(C.MmsValue_getType(mmsValue))
	return c.toGoValue(mmsValue, mmsType), nil
}

// ReadDataSet 读取DataSet
func (c *Client) ReadDataSet(objectRef string) ([]*MmsValue, error) {
	cObjectRef := C.CString(objectRef)
	defer C.free(unsafe.Pointer(cObjectRef))

	var clientError C.IedClientError
	dataSet := C.IedConnection_readDataSetValues(c.conn, &clientError, cObjectRef, nil)
	if err := GetIedClientError(clientError); err != nil {
		return nil, err
	}
	defer C.ClientDataSet_destroy(dataSet)

	dataSetValues := C.ClientDataSet_getValues(dataSet)
	// 长度
	dataSetSize := int(C.ClientDataSet_getDataSetSize(dataSet))
	mmsValues := make([]*MmsValue, dataSetSize)
	for i := 0; i < dataSetSize; i++ {
		value := C.MmsValue_getElement(dataSetValues, C.int(i))
		mmsType := MmsType(C.MmsValue_getType(value))
		mmsValue := &MmsValue{
			Type:  mmsType,
			Value: c.toGoValue(value, mmsType),
		}
		mmsValues[i] = mmsValue
	}
	return mmsValues, nil
}

// Close 关闭连接
func (c *Client) Close() {
	if c.conn != nil && c.connected.CompareAndSwap(true, false) {
		C.IedConnection_destroy(c.conn)
	}
}

// GetVariableSpecType 获取类型规格
func (c *Client) GetVariableSpecType(objectReference string, fc FC) (MmsType, error) {
	var clientError C.IedClientError
	cObjectRef := C.CString(objectReference)
	defer C.free(unsafe.Pointer(cObjectRef))

	// 获取类型
	spec := C.IedConnection_getVariableSpecification(c.conn, &clientError, cObjectRef, C.FunctionalConstraint(fc))
	if err := GetIedClientError(clientError); err != nil {
		return 0, err
	}
	defer C.MmsVariableSpecification_destroy(spec)
	mmsType := MmsType(C.MmsVariableSpecification_getType(spec))
	switch mmsType {
	case Integer:
		i := int(spec.typeSpec[0])
		switch i {
		case 8:
			return Int8, nil
		case 16:
			return Int16, nil
		case 32:
			return Int32, nil
		default:
			return Int64, nil
		}
	case Unsigned:
		switch int(spec.typeSpec[0]) {
		case 8:
			return Uint8, nil
		case 16:
			return Uint16, nil
		default:
			return Uint32, nil
		}
	}
	return mmsType, nil
}

func (c *Client) toGoValue(mmsValue *C.MmsValue, mmsType MmsType) interface{} {
	var value interface{}

	switch mmsType {
	case Integer:
		value = int64(C.MmsValue_toInt64(mmsValue))
	case Unsigned:
		value = uint32(C.MmsValue_toUint32(mmsValue))
	case Boolean:
		value = bool(C.MmsValue_getBoolean(mmsValue))
	case Float:
		value = float32(C.MmsValue_toFloat(mmsValue))
	case String, VisibleString:
		value = C.GoString(C.MmsValue_toString(mmsValue))
	case Structure, Array:
		value = c.toGoStructure(mmsValue, mmsType)
	case BitString:
		value = uint32(C.MmsValue_getBitStringAsInteger(mmsValue))
	case OctetString:
		size := uint16(C.MmsValue_getOctetStringSize(mmsValue))
		bytes := make([]byte, size)
		for i := 0; i < int(size); i++ {
			bytes[i] = uint8(C.MmsValue_getOctetStringOctet(mmsValue, C.int(i)))
		}
		value = bytes
	case BinaryTime:
		value = uint64(C.MmsValue_getBinaryTimeAsUtcMs(mmsValue))
	case UTCTime:
		value = uint32(C.MmsValue_toUnixTimestamp(mmsValue))
	}
	return value
}

func (c *Client) toGoStructure(mmsValue *C.MmsValue, mmsType MmsType) []*MmsValue {
	if mmsType != Structure {
		return nil
	}
	mmsValues := make([]*MmsValue, 0)
	for i := 0; ; i++ {
		value := C.MmsValue_getElement(mmsValue, C.int(i))
		// 读不到表示节点下没有属性了
		if value == nil {
			return mmsValues
		}
		valueType := MmsType(C.MmsValue_getType(value))
		mmsValues = append(mmsValues, &MmsValue{
			Type:  valueType,
			Value: c.toGoValue(value, valueType),
		})
	}
}

func (c *Client) getSubElementValue(sgcbVal *C.MmsValue, sgcbVarSpec *C.MmsVariableSpecification, name string) interface{} {
	mmsPath := C.CString(name)
	defer C.free(unsafe.Pointer(mmsPath))
	mmsValue := C.MmsValue_getSubElement(sgcbVal, sgcbVarSpec, mmsPath)
	defer C.MmsValue_delete(mmsValue)
	return c.toGoValue(mmsValue, MmsType(C.MmsValue_getType(mmsValue)))
}

// connect 建立连接
func connect(settings *Settings) (C.IedConnection, C.IedClientError) {
	conn := C.IedConnection_create()
	C.IedConnection_setConnectTimeout(conn, C.uint(settings.ConnectTimeout))
	C.IedConnection_setRequestTimeout(conn, C.uint(settings.RequestTimeout))
	host := C.CString(settings.Host)
	// 释放内存
	defer C.free(unsafe.Pointer(host))
	var err C.IedClientError
	C.IedConnection_connect(conn, &err, host, C.int(settings.Port))
	return conn, err
}
