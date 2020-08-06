package wjson

import (
	"fmt"
	"testing"
)

func TestMarshal(t *testing.T) {
	//obj:=[3]int{1,2}
	//obj := [3]string{"michael", "kobe", "iverson"}
	//var obj []int
	//obj = append(obj, 1, 2, 3)
	//obj := Student{
	//	Name:   "michael.w",
	//	ID:     1024,
	//	Class:  Class{3, 6},
	//	Pubkey: [5]int{1, 2, 3, 4, 5}}
	//typ := reflect.ValueOf(obj).Type()
	//field := typ.Field(0)
	//fmt.Println(getTag(string(field.Tag), "json"))
	//fmt.Println(len(field.Tag))
	obj := Student{"michael.w", 1024, Class{10, 20}, []int{1, 2, 3, 4, 5}}
	jsBytes, err := Marshal(obj)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(jsBytes))
	fmt.Println(len(jsBytes))

	var o Student
	if err = Unmarshal(jsBytes, &o); err != nil {
		panic(err)
	}
	fmt.Println(o)
}

type Student struct {
	Name   string `json:"name"`
	ID     uint
	Class  Class
	Pubkey []int
}

type Class struct {
	Grade int64
	Class int
}
