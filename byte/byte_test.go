// @program: json
// @author: edte
// @create: 2020-07-23 16:48
package byte

import (
	"testing"
)

func BenchmarkStringToBytes(b *testing.B) {
	b.ResetTimer()

	s := "soidjfopsajdfpsdofjisjdofjsodjfosdjiojsdopfjsdpoifjsopdjf"
	n := 1000000

	for i := 0; i < n; i++ {
		StringToBytes(s)
	}
}

func BenchmarkStringToBytesUnSafe(b *testing.B) {
	b.ResetTimer()

	s := "soidjfopsajdfpsdofjisjdofjsodjfosdjiojsdopfjsdpoifjsopdjf"
	n := 1000000

	for i := 0; i < n; i++ {
		StringToBytesUnSafe(s)
	}
}

func BenchmarkBytesToString(b *testing.B) {
	b.ResetTimer()
	by := []byte("sdofjsodfijosdjfosijdofjsodgjosdjgosjdgo")
	n := 10000000

	for i := 0; i < n; i++ {
		BytesToString(by)
	}
}

func BenchmarkBytesToStringUnSafe(b *testing.B) {
	b.ResetTimer()
	by := []byte("sdofjsodfijosdjfosijdofjsodgjosdjgosjdgo")
	n := 10000000

	for i := 0; i < n; i++ {
		BytesToStringUnSafe(by)
	}
}
