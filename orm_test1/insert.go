package orm_test1

import (
	"context"
	"orm_test1/internal/errs"
	"orm_test1/model"
)

type Inserter[T any] struct {
	builder
	sess           Session
	core           core
	values         []*T
	cols           []string
	onDuplicateKey *OnDuplicateKey
}

func (i Inserter[T]) Exec(ctx context.Context) Result {
	q, err := i.Build()
	if err != nil {
		return Result{
			err: err,
		}
	}
	res, err := i.sess.execContext(ctx, q.Sql, q.Args)
	return Result{
		res: res,
		err: err,
	}
}

func (i Inserter[T]) Build() (*Query, error) {
	i.sb.WriteString("INSERT INTO ")
	if len(i.values) == 0 {
		return nil, errs.ErrInsertZeroRow
	}
	m, err := i.core.r.Get(i.values[0])
	i.model = m
	if err != nil {
		return nil, err
	}
	i.builder.quota(m.TableName)
	i.sb.WriteString(" (")
	fields := m.Columns
	if len(i.cols) > 0 {
		fields = make([]*model.Field, 0, len(m.Columns))
		for _, col := range i.cols {
			fd, ok := i.model.FieldMap[col]
			if !ok {
				return nil, errs.NewErrField(col)
			}
			fields = append(fields, fd)
		}
	}
	for idx, fd := range fields {
		if idx > 0 {
			i.sb.WriteByte(',')
		}
		i.builder.quota(fd.Colname)
	}
	i.sb.WriteString(") ")
	i.sb.WriteString(" VALUES ")

	for idx, val := range i.values {
		if idx > 0 {
			i.sb.WriteByte(',')
		}
		i.sb.WriteByte('(')
		refv := i.core.creator(i.model, val)
		for jdx, fd := range fields {
			if jdx > 0 {
				i.sb.WriteByte(',')
			}
			i.sb.WriteByte('?')
			fdarg, err := refv.Field(fd.GoName)
			if err != nil {
				return nil, errs.NewErrUnknownColumn(fd.GoName)
			}
			i.args = append(i.args, fdarg)

		}
		i.sb.WriteByte(')')
	}
	//构造
	if i.onDuplicateKey != nil {
		err := i.dialect.buildDuplicateKey(&i.builder, i.onDuplicateKey)
		if err != nil {
			return nil, err
		}
	}

	i.sb.WriteByte(';')
	return &Query{
		Sql:  i.sb.String(),
		Args: i.args,
	}, nil
}

func (i *Inserter[T]) Columns(cols ...string) *Inserter[T] {
	i.cols = cols
	return i

}
func (i *Inserter[T]) Values(vals ...*T) *Inserter[T] {
	i.values = vals
	return i
}

func (i *Inserter[T]) OnDuplicatekey() *OnDuplicateKeyBuilder[T] {
	return &OnDuplicateKeyBuilder[T]{
		i: i,
	}
}
func NewInserter[T any](sess Session) *Inserter[T] {
	c := sess.getCore()
	return &Inserter[T]{
		core: c,
		sess: sess,
		builder: builder{
			core:    c,
			dialect: c.dialect,
		},
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
