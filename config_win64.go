//go:build windows && amd64

package iec61850

// #cgo CFLAGS: -I./libiec61850/windows_amd64/include
// #cgo LDFLAGS: -static-libgcc -static-libstdc++ -L${SRCDIR}/libiec61850/windows_amd64/lib -liec61850 -lws2_32 -liphlpapi
import "C"
