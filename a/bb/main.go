package main

import (
	"fmt"
	"unsafe"
)

type x struct {
	a bool
	b int16
	c []int
}

func main() {
	var x1 x
	fmt.Println(unsafe.Sizeof(x1))
	fmt.Println(unsafe.Alignof(x1))
}
