package internal

import (
	"errors"
	"reflect"
	"sync"
	"zorm/internal/model"
)

type Register interface {
	Get(typ reflect.Type) (*model.Model, error)
}

type DbRegistry struct {
	Models sync.Map
}

func (d *DbRegistry) Get(val any) (*model.Model, error) {
	typ := reflect.TypeOf(val)
	if typ.Kind() != reflect.Pointer || typ.Elem().Kind() != reflect.Struct {
		return nil,errors.New("")
	}
	model,ok := d.Models.Load(typ)
	if !ok  {

	}
}
func (d *DbRegistry)
