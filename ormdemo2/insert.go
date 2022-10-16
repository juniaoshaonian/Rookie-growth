package orm

import (
	"orm/internal/errs"
	model2 "orm/internal/model"
	"reflect"
)

type Inserter[T any] struct {
	builder
	values         []*T
	db             *DB
	cs             []string
	onDuplicateKey *OnDuplicateKey
}

func (i *Inserter[T]) Columns(cs ...string) *Inserter[T] {
	i.cs = cs
	return i
}

func (i *Inserter[T]) Build() (*Query, error) {
	if len(i.values) == 0 {
		return nil, errs.ErrInsertZeroRow
	}

	i.sb.WriteString("INSERT INTO ")
	m, err := i.db.r.Get(i.values[0])
	i.model = m
	if err != nil {
		return nil, err
	}
	i.sb.WriteByte('`')
	i.sb.WriteString(m.Tablename)
	i.sb.WriteByte('`')
	i.sb.WriteString(" (")
	fields := m.Cols
	if len(i.cs) != 0 {
		fields = make([]*model2.Field, 0, len(i.cs))
		for _, val := range i.cs {
			cval, ok := m.Fieldmap[val]
			if !ok {
				return nil, errs.NewErrUnknownField(val)
			}
			fields = append(fields, cval)
		}
	}

	for idx, val := range fields {
		if idx > 0 {
			i.sb.WriteByte(',')
		}
		i.sb.WriteByte('`')
		i.sb.WriteString(val.Colname)
		i.sb.WriteByte('`')
	}

	i.sb.WriteByte(')')
	i.sb.WriteString(" VALUES ")
	args := make([]any, 0, len(i.values)*len(m.Cols))
	for idx, val := range i.values {
		v := reflect.ValueOf(val).Elem()
		if idx > 0 {
			i.sb.WriteByte(',')
		}
		i.sb.WriteByte('(')

		for jdx, cval := range fields {
			if jdx > 0 {
				i.sb.WriteByte(',')
			}
			i.sb.WriteByte('?')
			fdVal := v.FieldByIndex(cval.Index)
			args = append(args, fdVal.Interface())
		}
		i.sb.WriteByte(')')
	}
	if i.onDuplicateKey != nil {
		err = i.dialect.buildDuplicateKey(&i.builder, i.onDuplicateKey)
		if err != nil {
			return nil, err
		}
	}
	i.sb.WriteByte(';')
	return &Query{
		SQL:  i.sb.String(),
		Args: args,
	}, nil
}

func NewInserter[T any](db *DB) *Inserter[T] {
	return &Inserter[T]{
		builder: builder{
			dialect: db.dialect,
		},
		db: db,
	}
}

func (i *Inserter[T]) Values(vals ...*T) *Inserter[T] {
	i.values = vals
	return i
}

func (i *Inserter[T]) OnDuplicateKey() *OnDuplicateKeyBuilder[T] {
	return &OnDuplicateKeyBuilder[T]{
		i: i,
	}
}

type OnDuplicateKeyBuilder[T any] struct {
	i *Inserter[T]
}

func (o *OnDuplicateKeyBuilder[T]) Update(assigns ...Assignable) *Inserter[T] {
	o.i.onDuplicateKey = &OnDuplicateKey{
		assigns: assigns,
	}
	return o.i
}

type OnDuplicateKey struct {
	assigns []Assignable
}
