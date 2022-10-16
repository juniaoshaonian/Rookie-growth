package model

import (
	"orm/internal/errs"
	"reflect"
	"sync"
	"unicode"
)

type Registery interface {
	Get(val any) (*Model, error)
	Register(val any, opt ...ModelOpt) (*Model, error)
}

type registery struct {
	Models sync.Map
}

func (r *registery) Get(val any) (*Model, error) {
	typ := reflect.TypeOf(val)
	m, ok := r.Models.Load(typ)
	if ok {
		return m.(*Model), nil
	}
	var err error
	m, err = r.Register(val)
	if err != nil {
		return nil, err
	}
	r.Models.Store(typ, m)
	return m.(*Model), nil
}

func (r *registery) Register(val any, opts ...ModelOpt) (*Model, error) {
	typ := reflect.TypeOf(val)
	if typ.Kind() != reflect.Ptr || typ.Elem().Kind() != reflect.Struct {
		return nil, errs.ErrPointerOnly
	}
	typ = typ.Elem()
	fieldNum := typ.NumField()
	fieldMap := make(map[string]*Field)
	colMap := make(map[string]*Field)
	cols := make([]*Field, 0, fieldNum)
	for i := 0; i < fieldNum; i++ {
		fd := typ.Field(i)
		colName := underscoreName(fd.Name)
		f := &Field{
			Typ:     fd.Type,
			Colname: colName,
			Goname:  fd.Name,
			Offset:  fd.Offset,
			Index:   fd.Index,
		}
		fieldMap[fd.Name] = f
		colMap[colName] = f
		cols = append(cols, f)
	}
	res := &Model{
		Tablename: underscoreName(typ.Name()),
		Fieldmap:  fieldMap,
		Colmap:    colMap,
		Cols:      cols,
	}
	for _, opt := range opts {
		opt(res)
	}
	return res, nil

}
func underscoreName(tableName string) string {
	var buf []byte
	for i, v := range tableName {
		if unicode.IsUpper(v) {
			if i != 0 {
				buf = append(buf, '_')
			}
			buf = append(buf, byte(unicode.ToLower(v)))
		} else {
			buf = append(buf, byte(v))
		}

	}
	return string(buf)
}
func NewRegitery() Registery {
	return &registery{}
}
