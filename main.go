package main

import (
	"fmt"

	"github.com/paniccaaa/wbtech/develop/dev02"
)

func main() {
	str, err := dev02.Unpack("45")
	if err != nil {
		panic(err)
	}
	fmt.Println(str, str == "aaaabccddddde") // my = aaaabcccdddddde || expected = aaaabccddddde
}
