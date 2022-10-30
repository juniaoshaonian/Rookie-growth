package model

import (
	"orm_test1/internal/errs"
	"reflect"
	"sync"
	"unicode"
)

type Option func(model *Model) error
type Registry interface {
	Get(val any) (*Model, error)
	Register(val any, opts ...Option) (*Model, error)
}

type Register struct {
	Models sync.Map
}

func (r *Register) Get(val any) (*Model, error) {
	t := reflect.TypeOf(val)
	v, ok := r.Models.Load(t)
	if ok {
		return v.(*Model), nil
	}
	return r.Register(val)
}

func (r *Register) Register(val any, opts ...Option) (*Model, error) {
	m, err := r.ParseModel(val)
	if err != nil {
		return nil, err
	}
	for _, opt := range opts {
		err = opt(m)
		if err != nil {
			return nil, err
		}
	}
	r.Models.Store(reflect.TypeOf(val), m)
	return m, nil
}
func (r *Register) ParseModel(val any) (*Model, error) {
	//只允许注册指针
	if val == nil {
		return nil, errs.ErrPointerOnly
	}
	reflectVal := reflect.ValueOf(val)
	if reflectVal.Kind() != reflect.Ptr || reflectVal.Elem().Kind() != reflect.Struct {
		return nil, errs.ErrPointerOnly
	}
	reflectVal = reflectVal.Elem()
	reflectType := reflectVal.Type()
	tableName := underscoreName(reflectType.Name())
	fieldMap := make(map[string]*Field, reflectType.NumField())
	ColMap := make(map[string]*Field, reflectType.NumField())
	Colunms := make([]*Field, 0, reflectType.NumField())
	for i := 0; i < reflectType.NumField(); i++ {
		field := &Field{
			Colname: underscoreName(reflectType.Field(i).Name),
			Type:    reflectType.Field(i).Type,
			GoName:  reflectType.Field(i).Name,
			Index:   reflectType.Field(i).Index,
			Offset:  reflectType.Field(i).Offset,
		}
		fieldMap[field.GoName] = field
		ColMap[field.Colname] = field
		Colunms = append(Colunms, field)
	}
	model := &Model{
		TableName: tableName,
		FieldMap:  fieldMap,
		Columns:   Colunms,
		ColumnMap: ColMap,
	}
	r.Models.Store(tableName, model)
	return model, nil

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

func NewRegistry() Registry {
	return &Register{}
}
