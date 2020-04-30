package main

import (
	"fmt"
	"github.com/vfw4g/base/utils/conv"
)

type Foo struct {
	Name  string
	Age   int
	Roles []string
}

type Bar struct {
	Age   int
	Name  string `cp:"name"`
	Roles []string
}

func main() {
	f := Foo{
		Name: "kook",
		Age:  18,
		Roles: []string{
			"admin",
		},
	}
	b := Bar{}

	if err := conv.StructFieldCopy(&f, &b); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", b)
}
