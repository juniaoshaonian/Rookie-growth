package model

import (
	"reflect"
)

type Model struct {
	Tablename string
	Fieldmap  map[string]*Field
	Colmap    map[string]*Field
	Cols      []*Field
}
type Field struct {
	Colname string
	Goname  string
	Typ     reflect.Type
	Offset  uintptr
	Index   []int
}

type ModelOpt func(model *Model)

func WithTableName(name string) ModelOpt {
	return func(model *Model) {
		model.Tablename = name
	}
}

func WithColunmName(key string, name string) ModelOpt {
	return func(model *Model) {
		model.Fieldmap[key].Colname = name
	}
}
