//go:build windows && amd64

package iec61850

// #cgo CFLAGS: -I./libiec61850/win64/include
// #cgo LDFLAGS: -static-libgcc -static-libstdc++ -Wl,--start-group ${SRCDIR}/libiec61850/win64/lib/libiec61850.a ${SRCDIR}/libiec61850/win64/lib/libhal.a ${SRCDIR}/libiec61850/win64/lib/libwpcap.a ${SRCDIR}/libiec61850/win64/lib/libpacket.a -Wl,--end-group -lws2_32 -liphlpapi
import "C"
