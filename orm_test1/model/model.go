package model

import "reflect"

type Model struct {
	TableName string
	FieldMap  map[string]*Field
	ColumnMap map[string]*Field
	Columns   []*Field
}

type Field struct {
	Colname string
	GoName  string
	Type    reflect.Type
	Index   []int
	Offset  uintptr
}
