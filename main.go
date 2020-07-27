// @program: json
// @author: edte
// @create: 2020-07-18 19:58
package main

import (
	"fmt"
	"log"

	"json/encode"
)

func main() {
	a := []int{23, 4, 5, 2}

	marshal, err := encode.Marshal(&a)
	if err != nil {
		log.Printf("failed to : %v\n", err)
		return
	}
	fmt.Println(string(marshal))
}
