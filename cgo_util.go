package iec61850

import "C"
import sc "golang.org/x/text/encoding/simplifiedchinese"

func C2GoStr(str *C.char) string {
	utf8str, _ := sc.GB18030.NewDecoder().String(C.GoString(str))
	return utf8str
}

func Go2CStr(str string) *C.char {
	gbstr, _ := sc.GB18030.NewEncoder().String(str)
	return C.CString(gbstr)
}

func C2GoBool(i C.int) bool {
	if i == 1 {
		return true
	}
	return false
}

func Go2CBool(b bool) C.int {
	if b {
		return 1
	}
	return 0
}

func IsBitSet(n int, pos uint) bool {
	mask := 1 << pos
	return (n & mask) != 0
}
