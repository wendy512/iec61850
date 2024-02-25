//go:build linux && arm64

package iec61850

// #cgo CFLAGS: -I./libiec61850/inc/hal/inc -I./libiec61850/inc/common/inc -I./libiec61850/inc/goose -I./libiec61850/inc/iec61850/inc -I./libiec61850/inc/iec61850/inc_private -I./libiec61850/inc/logging -I./libiec61850/inc/mms/inc -I./libiec61850/inc/mms/inc_private -I./libiec61850/inc/mms/iso_mms/asn1c
// #cgo LDFLAGS: -static-libgcc -static-libstdc++ -L./libiec61850/lib/linux_armv8 -liec61850 -lhal
import "C"
