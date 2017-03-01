package ref

import (
	"fmt"
	"reflect"
	"test/reflect/testmodule"
)

type MyStruct struct {
	name string
}

func (this *MyStruct) GetName(str string) string {
	this.name = str
	return this.name
}

func TestReflect1() {
	s := "this is a string"
	fmt.Println("reflect typeof :", reflect.TypeOf(s))

	fmt.Println("reflect valueof:", reflect.ValueOf(s))

	fmt.Println("=====================")
	var x float64 = 3.4
	fmt.Println("reflect typeof :", reflect.TypeOf(x))

	fmt.Println("reflect valueof:", reflect.ValueOf(x))

	fmt.Println("=====================")

	a := MyStruct{name: "tac"}

	fmt.Println("struct reflect :", reflect.TypeOf(a))
	// fmt.Println("struct reflect :", reflect.TypeOf(a).Method(0).Name)

	fmt.Println("print type method")

	// fmt.Println(reflect.TypeOf(a).NumField())
	// fmt.Println(reflect.TypeOf(a).Field(0).Name)
	fmt.Println(reflect.TypeOf(a).NumMethod())
	fmt.Println("=====================")
	fmt.Println(reflect.ValueOf(a).NumMethod())
	fmt.Println("=====================")
	fmt.Println(reflect.ValueOf(a).NumField()) //NumField 必须是一个结构体类型或者具体类型的实例，不能是一个指针

	fmt.Println("==========into for loop=========")
	for m := 0; m < reflect.TypeOf(a).NumMethod(); m++ {
		method := reflect.TypeOf(a).Method(m)
		fmt.Println("type:", method.Type)
		fmt.Println("name:", method.Name)
		fmt.Println("params num:", method.Type.NumIn())
		fmt.Println("params type:", method.Type.In(1))
	}

	fmt.Println("=====================")
	fmt.Println(reflect.ValueOf(a))

	fmt.Println("=====================")
	fmt.Println(reflect.Indirect(reflect.ValueOf(a)).Type().Name())
	fmt.Println("=====================")
	fmt.Println(reflect.ValueOf(a).Type().Name())
}

func TestReflect2() {
	m := &module.Module{User: "Gavin", Age: 26, Sex: "male"}
	DoFunction(m)
}

func DoFunction(obj interface{}) {
	fmt.Println("in=======")
	for i := 0; i < reflect.TypeOf(obj).NumMethod(); i++ { //如果obj是一个结构体实例的地址，那么可以通过NumMethod获取这个类型的所有绑定方法。如果只是传入实例，也就是结构体，那么无法获取方法个数
		if reflect.TypeOf(obj).Method(i).Name == "GetUser" {
			fmt.Println("this is function getuser")
			continue
		}
		name := "tac"
		in := make([]reflect.Value, 1)
		in[0] = reflect.ValueOf(name)
		ret := reflect.ValueOf(obj).Method(i).Call(in)
		iter := ret[0].Interface()
		if value, ok := iter.(int); ok {
			fmt.Println("value:", value)
		} else {
			value := iter.(string)
			fmt.Println("value:", value)
		}
	}
	fmt.Println("out======")
}
