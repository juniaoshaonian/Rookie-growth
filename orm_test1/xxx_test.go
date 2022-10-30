package orm_test1

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_XX(t *testing.T) {
	refU := reflect.ValueOf(&V2{})
	indx := map[string][]int{}
	offsets := map[string]uintptr{}
	refU = refU.Elem()
	Type(refU.Type(), indx, offsets)
	fmt.Println("a")
	xx := refU.FieldByIndex([]int{0, 0, 0}).Type().Name()
	fmt.Println(xx)
}
func Type(typ reflect.Type, indx map[string][]int, offsets map[string]uintptr) {
	for i := 0; i < typ.NumField(); i++ {
		fd := typ.Field(i)
		if fd.Type.Kind() == reflect.Struct {
			Type(fd.Type, indx, offsets)
		} else {
			indx[fd.Name] = fd.Index
			offsets[fd.Name] = fd.Offset
		}
	}
}

type User struct {
	Id int
	BaseEntity
	Age  int8
	Name string
}

type BaseEntity struct {
	CreateTime int64
	UpdateTime int64
}

type V1 struct {
	BaseEntity
}
type V2 struct {
	V1
}
