// @program: json
// @author: edte
// @create: 2020-07-18 19:58
package main

import (
	"fmt"
	"log"

	"json/value"
)

/*func main() {
	j := []byte(`{"name":34,"fs":"gsd","ha":{"fsd":23,"ow":3}}`)

	var i interface{}

	json.Unmarshal(j, &i)

	m := i.(map[string]interface{})
	fmt.Println(m) // map[fs:gsd ha:map[fsd:23 ow:3] name:34]

	bb := m["ha"].(map[string]interface{})
	bb["fsd"] = "ah"
	fmt.Println(bb) // map[fsd:23 ow:3]
	m["ha"] = bb

	fmt.Println(m)
}
*/

func main() {
	d := []byte(`{"name":34,"fs":"gsd","ha":{"fsd":{"hhh":2},"ow":3}}`)

	j, err := value.NewJson(d)
	if err != nil {
		log.Printf("failed to new json:%s", err)
		return
	}

	var m interface{}
	m = `{"sdf":34,"o":[1,2,"sdf"]}`

	err = j.Set(m, []string{"ha", "fsd", "hhh"})
	if err != nil {
		log.Printf("failed to set json value:%s", err)
		return
	}
	marshal, err := j.Marshal()
	fmt.Println(marshal)
}
