package iec61850

// #include <tls_config.h>
import "C"
import (
	"fmt"
	"unsafe"
)

type TLSConfigVersion int

const (
	TLS_VERSION_NOT_SELECTED TLSConfigVersion = 0
	TLS_VERSION_SSL_3_0      TLSConfigVersion = 3
	TLS_VERSION_TLS_1_0                       = 4
	TLS_VERSION_TLS_1_1                       = 5
	TLS_VERSION_TLS_1_2                       = 6
	TLS_VERSION_TLS_1_3                       = 7
)

type TLSConfigurationEventHandler func(parameter unsafe.Pointer, eventLevel, eventCode int, message string, conn C.TLSConnection)

type TLSConfig struct {
	KeyFile                      string // Path to the key file
	KeyPassword                  string // Password for the key file
	CertFile                     string // Path to the certificate file
	ChainValidation              bool   // Enable chain validation
	AllowOnlyKnownCertificates   bool   // Allow only known certificates
	MinTlsVersion                TLSConfigVersion
	MaxTlsVersion                TLSConfigVersion
	caCerts                      []string
	allowedCertificates          []string
	tlsConfigurationEventHandler *TLSConfigurationEventHandler
}

func NewTLSConfig() *TLSConfig {
	return &TLSConfig{
		ChainValidation:            true,
		AllowOnlyKnownCertificates: false,
		MinTlsVersion:              TLS_VERSION_TLS_1_0,
		MaxTlsVersion:              TLS_VERSION_NOT_SELECTED,
		caCerts:                    make([]string, 0),
		allowedCertificates:        make([]string, 0),
	}
}

func (that *TLSConfig) AddCACertificateFromFile(filename string) {
	that.caCerts = append(that.caCerts, filename)
}

func (that *TLSConfig) AddAllowedCertificateFromFile(filename string) {
	that.allowedCertificates = append(that.allowedCertificates, filename)
}

func (that *TLSConfig) SetEventHandler(handler *TLSConfigurationEventHandler) {
	that.tlsConfigurationEventHandler = handler
}

func (that *TLSConfig) createCTlsConfig() (C.TLSConfiguration, error) {
	cKeyFile := C.CString(that.KeyFile)
	defer C.free(unsafe.Pointer(cKeyFile))

	cKeyPassword := C.CString(that.KeyPassword)
	defer C.free(unsafe.Pointer(cKeyPassword))

	cCertFile := C.CString(that.CertFile)
	defer C.free(unsafe.Pointer(cCertFile))

	tlsConfig := C.TLSConfiguration_create()
	C.TLSConfiguration_setChainValidation(tlsConfig, C.bool(that.ChainValidation))
	C.TLSConfiguration_setAllowOnlyKnownCertificates(tlsConfig, C.bool(that.AllowOnlyKnownCertificates))
	C.TLSConfiguration_setMinTlsVersion(tlsConfig, C.TLSConfigVersion(that.MinTlsVersion))
	C.TLSConfiguration_setMaxTlsVersion(tlsConfig, C.TLSConfigVersion(that.MaxTlsVersion))

	if that.KeyPassword == "" {
		if !bool(C.TLSConfiguration_setOwnKeyFromFile(tlsConfig, cKeyFile, nil)) {
			return nil, fmt.Errorf("failed to load private key %s", that.KeyFile)
		}
	} else {
		if !bool(C.TLSConfiguration_setOwnKeyFromFile(tlsConfig, cKeyFile, cKeyPassword)) {
			return nil, fmt.Errorf("failed to load private key %s", that.KeyFile)
		}
	}

	if !bool(C.TLSConfiguration_setOwnCertificateFromFile(tlsConfig, cCertFile)) {
		return nil, fmt.Errorf("failed to load own certificate %s", that.CertFile)
	}

	for _, caCert := range that.caCerts {
		cCACert := C.CString(caCert)
		if !bool(C.TLSConfiguration_addCACertificateFromFile(tlsConfig, cCACert)) {
			C.free(unsafe.Pointer(cCACert))
			return nil, fmt.Errorf("failed to load CA certificate %s", caCert)
		}
		C.free(unsafe.Pointer(cCACert))
	}

	for _, cert := range that.allowedCertificates {
		cCert := C.CString(cert)
		if !bool(C.TLSConfiguration_addAllowedCertificateFromFile(tlsConfig, cCert)) {
			C.free(unsafe.Pointer(cCert))
			return nil, fmt.Errorf("failed to load allowed certificate %s", cert)
		}
		C.free(unsafe.Pointer(cCert))
	}

	return tlsConfig, nil
}
