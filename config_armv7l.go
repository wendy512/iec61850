//go:build linux && arm

package iec61850

// #cgo CFLAGS: -I./libiec61850/linux_armv7/include
// #cgo LDFLAGS: -static-libgcc -static-libstdc++ -L./libiec61850/linux_armv7/lib -liec61850 -lpthread
import "C"
