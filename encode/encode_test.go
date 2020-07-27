// @program: json
// @author: edte
// @create: 2020-07-27 21:33
package encode

import (
	"encoding/json"
	"testing"
)

func TestMarshal(t *testing.T) {
	var datas = []interface{}{
		// struct
		struct {
			Age  int
			Name string
			MM   interface{}
		}{
			Age:  18,
			Name: "lily",
			MM:   "sdf",
		},
		// map
		map[string]interface{}{"sfd": 234, "objso": "sdf", "sdf": false, "sdfoj": map[string]int{"sdf": 234,
			"oj": 234}},
		// bool
		false,
		true,
		// number
		1,
		1.234,
		// string
		"sdfi",
		// slice
		[]int{1, 3, 45},
		[]string{"f", "sdf2", "234"},
		[]interface{}{23, "234", false},
		// array
		[3]int{1, 3, 45},
		[5]string{"f", "sdf2", "234"},
		[9]interface{}{23, "234", false},
		// ptr
		&[]int{234, 5, 5},
	}

	for i, data := range datas {
		j, err := Marshal(data)
		if err != nil {
			t.Errorf("%d: get  %v, err: %v \n", i, string(j), err)
		}
	}
}

// orz, 竟然比标准库还慢 4 倍
func BenchmarkMarshal(b *testing.B) {
	j := struct {
		Age    int
		Name   string
		Butty  bool
		High   float64
		Family []int
		Cats   [5]string
		HH     map[string]int
		School struct {
			Name string
			Age  int
		}
	}{
		Age:    18,
		Name:   "lily",
		Butty:  false,
		High:   20.1,
		Family: []int{1, 4, 5},
		Cats:   [5]string{"fd", "wo"},
		HH:     map[string]int{"fso": 234, "235": 99},
		School: struct {
			Name string
			Age  int
		}{
			Name: "lucy",
			Age:  49,
		},
	}
	b.ResetTimer()
	for i := 0; i < 100000; i++ {
		Marshal(j)
	}
}

func BenchmarkMarshal2(b *testing.B) {
	j := struct {
		Age    int
		Name   string
		Butty  bool
		High   float64
		Family []int
		Cats   [5]string
		HH     map[string]int
		School struct {
			Name string
			Age  int
		}
	}{
		Age:    18,
		Name:   "lily",
		Butty:  false,
		High:   20.1,
		Family: []int{1, 4, 5},
		Cats:   [5]string{"fd", "wo"},
		HH:     map[string]int{"fso": 234, "235": 99},
		School: struct {
			Name string
			Age  int
		}{
			Name: "lucy",
			Age:  49,
		},
	}
	b.ResetTimer()
	for i := 0; i < 100000; i++ {
		json.Marshal(j)
	}
}
