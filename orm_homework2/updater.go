package orm

import (
	"context"
	"orm/internal/errs"
	"reflect"
)

type Updater[T any] struct {
	builder
	db      *DB
	assigns []Assignable
	val     *T
	where   []Predicate
}

func NewUpdater[T any](db *DB) *Updater[T] {
	return &Updater[T]{
		db: db,
		builder: builder{
			dialect: db.dialect,
			quoter:  db.dialect.quoter(),
		},
	}
}

func (u *Updater[T]) Update(t *T) *Updater[T] {
	u.val = t
	return u
}

func (u *Updater[T]) Set(assigns ...Assignable) *Updater[T] {
	u.assigns = assigns
	return u
}

func (u *Updater[T]) Build() (*Query, error) {
	if len(u.assigns) == 0 {
		return nil, errs.ErrNoUpdatedColumns
	}
	t := new(T)
	model, err := u.db.r.Get(t)
	if err != nil {
		return nil, err
	}
	u.model = model
	u.sb.WriteString("UPDATE ")
	u.quote(model.TableName)
	u.sb.WriteString(" SET ")
	for idx, as := range u.assigns {
		if idx > 0 {
			u.sb.WriteByte(',')
		}
		switch expr := as.(type) {
		case Column:
			fd, ok := u.model.FieldMap[expr.name]
			if !ok {
				return nil, errs.NewErrUnknownField(expr.name)
			}
			u.quote(fd.ColName)
			u.sb.WriteString("=?")
			rVal := reflect.ValueOf(u.val).Elem()
			fdVal := rVal.Field(fd.Index)
			u.args = append(u.args, fdVal.Interface())
		case Assignment:
			fd, ok := u.model.FieldMap[expr.column]
			if !ok {
				return nil, errs.NewErrUnknownField(expr.column)
			}
			u.quote(fd.ColName)
			u.sb.WriteByte('=')
			u.buildExpression(expr.val)
		}
	}
	if len(u.where) > 0 {
		u.sb.WriteString(" WHERE ")
		if err = u.buildPredicates(u.where); err != nil {
			return nil, err
		}
	}
	u.sb.WriteByte(';')
	return &Query{
		SQL:  u.sb.String(),
		Args: u.args,
	}, nil

}

func (u *Updater[T]) Where(ps ...Predicate) *Updater[T] {
	u.where = ps
	return u
}

func (u *Updater[T]) Exec(ctx context.Context) Result {
	return Result{}
}

// AssignNotZeroColumns 更新非零值
func AssignNotZeroColumns(entity interface{}) []Assignable {
	entVal := reflect.ValueOf(entity)
	entType := reflect.TypeOf(entity)
	for entVal.Kind() == reflect.Ptr {
		entVal = entVal.Elem()
		entType = entType.Elem()
	}
	numField := entType.NumField()
	Assigns := make([]Assignable, 0, numField)
	for i := 0; i < entType.NumField(); i++ {
		if !entVal.Field(i).IsZero() {
			a := Assignment{
				column: entType.Field(i).Name,
				val: value{
					val: entVal.Field(i).Interface(),
				},
			}
			Assigns = append(Assigns, a)
		}
	}
	return Assigns
}
