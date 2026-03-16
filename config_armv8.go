//go:build linux && arm64

package iec61850

// #cgo CFLAGS: -I./libiec61850/linux_armv8/include
// #cgo LDFLAGS: -static-libgcc -static-libstdc++ -L./libiec61850/linux_armv8/lib -liec61850 -lpthread
import "C"
