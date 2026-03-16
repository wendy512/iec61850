//go:build linux && amd64

package iec61850

// #cgo CFLAGS: -I./libiec61850/linux_amd64/include
// #cgo LDFLAGS: -static-libgcc -static-libstdc++ -L./libiec61850/linux_amd64/lib -liec61850 -lpthread
import "C"
