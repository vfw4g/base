package main

import (
	"fmt"
	"strings"
)

type Student struct {
	Name string
}

func (s *Student) SayName() string {
	return s.Name
}

func initStudent() *Student {
	return &Student{
		Name: "susan",
	}
}

func main() {
	//var s = Student{
	//	Name: "kook",
	//}
	//fmt.Println(unsafe.Pointer(&s))
	//var s1 interface{}
	//s1 = &s //old
	//fmt.Println(unsafe.Pointer(&s1))
	//a := &s1
	//fmt.Println(unsafe.Pointer(a))
	////reflect.ValueOf(s1).
	//s2 := &Student{
	//	Name: "haha",
	//}
	//v2 := reflect.ValueOf(s2).Elem()
	//reflect.ValueOf(s1).Elem().Set(v2)
	//fmt.Println(s.SayName())
	//fmt.Println(*s1)
	//var s2 interface{}//new
	//a := s1.(*Student)
	//fmt.Println(a.SayName())
	//s2 := &Student{
	//	Name:"haha",
	//}
	//fmt.Println((*int16)(unsafe.Pointer(&s2)))
	//s3 := &s1
	//fmt.Println((*int16)(unsafe.Pointer(&s3)))
	//s4 := (*int16)(unsafe.Pointer(&s1))
	//fmt.Println(s4)
	//fmt.Println((*s3).(*Student).SayName())
	//s3 := s2
	//fmt.Println((*s3).(*Student).SayName())
	//addS2 := (*Student)(unsafe.Pointer(s2))
	//s1 = s2
	//b := addS2.(*Student)
	//fmt.Println(s1.SayName())
	ss := strings.Split("hello,", ",")
	fmt.Println(ss, len(ss))
}
