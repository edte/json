// @program: json
// @author: edte
// @create: 2020-07-23 23:35
package scanner

import (
	"testing"
)

func TestIsValid(t *testing.T) {
	tests := []struct {
		name  string
		datas []byte
		want  bool
	}{
		// {"number", []byte(`1`), true},
		// {"number", []byte(`1.`), true},
		// {"number", []byte(`1.2`), true},
		// {"number", []byte(`-1`), true},
		// {"number", []byte(`-1.`), true},
		// {"number", []byte(`-1.2`), true},
		// {"number", []byte(`+`), true},
		// {"float", []byte(`1.2`), true},
		// {"string", []byte(`"sdfsd"`), true},
		// {"false", []byte(`false`), true},
		// {"true", []byte(`true`), true},
		// {"null", []byte(`null`), true},
		// {"null", []byte(`anull`), true},
		// {"null", []byte(`nullb`), true},
		{"array", []byte(`[null]`), true},
		{"array", []byte(`[1]`), true},
		{"array", []byte(`[1,2]`), true},
		{"array", []byte(`[1,2,"sd"]`), true},
		{"array", []byte(`[1,false,"sd"]`), true},
		{"array", []byte(`[[2],false,"sd"]`), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValid(tt.datas); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
