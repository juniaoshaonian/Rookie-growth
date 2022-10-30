package valuer

import (
	"database/sql"
	"orm_test1/internal/errs"
	"orm_test1/model"
	"reflect"
	"unsafe"
)

type unsafevaluer struct {
	t     any
	model *model.Model
	addr  unsafe.Pointer
}

func (u *unsafevaluer) Field(name string) (any, error) {
	fd, ok := u.model.FieldMap[name]
	if !ok {
		return nil, errs.NewErrField(name)
	}
	ptr := unsafe.Pointer(uintptr(u.addr) + fd.Offset)
	val := reflect.NewAt(fd.Type, ptr).Elem()
	return val.Interface(), nil
}

func (u *unsafevaluer) SetFields(rows *sql.Rows) error {
	cols, err := rows.Columns()
	if err != nil {
		return err
	}
	relVal := make([]any, 0, len(cols))
	for _, col := range cols {
		fd, ok := u.model.ColumnMap[col]
		if !ok {
			return errs.NewErrUnknownColumn(col)
		}
		fdval := reflect.NewAt(fd.Type, unsafe.Pointer(uintptr(u.addr)+fd.Offset)).Interface()
		relVal = append(relVal, fdval)
	}
	rows.Scan(relVal...)
	return nil
}

func NewUnsafevaluer(m *model.Model, t any) Valuer {
	addr := unsafe.Pointer(reflect.ValueOf(t).Pointer())
	return &unsafevaluer{
		t:     t,
		addr:  addr,
		model: m,
	}
}
