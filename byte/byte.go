// @program: json
// @author: edte
// @create: 2020-07-23 16:20
package byte

import (
	"unsafe"
)

func BytesToStringUnSafe(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func StringToBytesUnSafe(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

func StringToBytes(s string) []byte {
	return []byte(s)
}

func BytesToString(b []byte) string {
	return string(b)
}
