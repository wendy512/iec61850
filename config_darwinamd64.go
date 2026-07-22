//go:build darwin && amd64

package iec61850

// #cgo CFLAGS: -I./libiec61850/darwin_amd64/include
// #cgo LDFLAGS: -L./libiec61850/darwin_amd64/lib -liec61850 -lpthread
import "C"
