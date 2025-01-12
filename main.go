package main

import (
	"github.com/paniccaaa/wbtech/develop/dev01"
)

func main() {
	_, err := dev01.Now()
	if err != nil {
		panic(err)
	}
}
