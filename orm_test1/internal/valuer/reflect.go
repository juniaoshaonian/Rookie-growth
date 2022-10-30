package valuer

import (
	"database/sql"
	"orm_test1/internal/errs"
	"orm_test1/model"
	"reflect"
)

type reflectValuer struct {
	model *model.Model
	t     reflect.Value
}

func (r *reflectValuer) Field(name string) (any, error) {
	val := r.t
	typ := val.Type()
	_, ok := typ.FieldByName(name)
	if !ok {
		return nil, errs.NewErrField(name)
	}

	return val.FieldByName(name).Interface(), nil
}

func (r *reflectValuer) SetFields(rows *sql.Rows) error {
	cols, err := rows.Columns()
	if err != nil {
		return err
	}
	refVal := make([]any, 0, len(cols))
	eles := make([]reflect.Value, 0, len(cols))
	for _, col := range cols {
		fd, ok := r.model.ColumnMap[col]
		if !ok {
			return errs.NewErrUnknownColumn(col)
		}
		xx := reflect.New(fd.Type)
		refVal = append(refVal, xx.Interface())
		eles = append(eles, xx.Elem())
	}
	rows.Scan(refVal...)
	refT := r.t
	for index, col := range cols {
		ft := refT.FieldByIndex(r.model.ColumnMap[col].Index)
		ft.Set(eles[index])
	}
	return nil
}
func NewReflectValuer(m *model.Model, t any) Valuer {
	return &reflectValuer{
		model: m,
		t:     reflect.ValueOf(t).Elem(),
	}
}
