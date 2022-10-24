package model

import "reflect"

type Model struct {
}

type Field struct {
	Colname string
	GoName  string
	Type    reflect.Type
}
