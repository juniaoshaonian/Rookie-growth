package valuer

import (
	"database/sql"
	"orm/internal/errs"
	"orm/internal/model"
	"reflect"
)

type reflectValue struct {
	t     any
	model *model.Model
}

func (r reflectValue) SetColumns(rows *sql.Rows) error {
	cols, err := rows.Columns()
	if err != nil {
		return err
	}
	if len(cols) > len(r.model.Fieldmap) {
		return errs.ErrTooManyReturnedColumns
	}

	res := make([]any, len(cols))
	eles := make([]reflect.Value, len(cols))
	for index, col := range cols {
		fd, ok := r.model.Colmap[col]
		if !ok {
			return errs.NewErrUnknownColumn(col)
		}
		res[index] = reflect.New(fd.Typ).Interface()
		eles[index] = reflect.ValueOf(res[index]).Elem()
	}
	rows.Scan(res...)
	t := r.t
	valt := reflect.ValueOf(t).Elem()
	for index, col := range cols {
		f := r.model.Colmap[col]
		tt := valt.FieldByName(f.Goname)
		tt.Set(eles[index])
	}
	return nil
}

func NewReflectValue(t any, m *model.Model) Valuer {
	return reflectValue{
		t:     t,
		model: m,
	}
}
