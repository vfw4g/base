package errors

import (
	"fmt"
	"reflect"
	"testing"
)

type hello interface {
}

func f1() func() {
	return nil
}

func TestCause(t *testing.T) {
	e1 := f1()
	fmt.Println(e1 == nil)
	//fmt.Println("1",unsafe.Pointer(&e1),)
	if err := reflect.ValueOf(e1).Interface(); err != nil {
		fmt.Println("2", err)
		fmt.Println("not nil")
	}
	fmt.Println("ok")
}
