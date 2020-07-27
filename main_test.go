// @program: json
// @author: edte
// @create: 2020-07-19 19:47
package main

import (
	"testing"

	"json/encode"
)

/*func BenchmarkIsValid(b *testing.B) {
	c := `{a"a":2}`
	b.ResetTimer()
	for i := 0; i < 10000; i++ {
		IsValid(c)
	}

}
*/
/*func BenchmarkIsValid2(b *testing.B) {
	c := `{a"a":2}`
	b.ResetTimer()
	for i := 0; i < 10000; i++ {
		Valid([]byte(c))
	}
}
*/

func BenchmarkA(b *testing.B) {
	c := `{"a":2}`
	b.ResetTimer()
	for i := 0; i < 1000000; i++ {
		validpayload(stringBytes(c), 0)
	}
}

func BenchmarkB(b *testing.B) {
	c := `{"a":2}`
	b.ResetTimer()
	for i := 0; i < 1000000; i++ {
		encode.Valid(stringBytes(c))
	}
}
