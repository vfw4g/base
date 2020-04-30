package bean

import (
	"fmt"
)

type Student struct {
	Name string
}

func (s *Student) SayName() string {
	return s.Name
}

func initStudent() Instance {
	return &Student{
		Name: "susan",
	}
}

func ExampleGetBeanByNameDelay() {
	s := Student{
		Name: "test",
	}
	GetBeanByNameDelay("student", &s)
	RegisterSingletonBean("student", initStudent)
	fmt.Println(s.SayName())
	//output:
	//susan
}

func ExampleGetBeanByName() {
	s := Student{
		Name: "test",
	}
	GetBeanByNameDelay("student", &s)
	RegisterSingletonBean("student", initStudent)
	s2 := GetBeanByName("student").(*Student)
	fmt.Println(s2.SayName())
	//output:
	//susan
}

func ExampleRegisterBeforeGet() {
	RegisterSingletonBean("student", initStudent)
	s2 := GetBeanByName("student").(*Student)
	fmt.Println(s2.SayName())
	//output:
	//susan
}

func ExampleDoubleRegister() {
	RegisterSingletonBean("student", initStudent)
	RegisterSingletonBean("student", initStudent)
	s2 := GetBeanByName("student").(*Student)
	fmt.Println(s2.SayName())
	//output:
	//susan
}
