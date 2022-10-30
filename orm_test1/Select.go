package orm_test1

import (
	"context"
	"errors"
	"orm_test1/internal/errs"
)

type Select[T any] struct {
	typ  string
	sess Session
	core
	table TableReference
	builder
	where []*Predicate
	cols  []Selectable
}

func (s *Select[T]) Get(ctx context.Context) (*T, error) {
	model, err := s.r.Get(new(T))
	if err != nil {
		return nil, err
	}
	res := get[T](ctx, s.core, s.sess, &QueryContext{
		Type:    s.typ,
		Builder: s,
		Model:   model,
	})
	if res.Res != nil {
		return res.Res.(*T), nil
	}
	return nil, res.Err
}

func (s *Select[T]) GetMulti(ctx context.Context) ([]*T, error) {
	q, err := s.Build()
	if err != nil {
		return nil, err
	}
	rows, err := s.sess.queryContext(ctx, q.Sql, q.Args)
	if err != nil {
		return nil, err
	}
	ans := make([]*T, 0, 16)
	for rows.Next() {
		t := new(T)
		valueGet := s.core.creator(s.model, t)
		err = valueGet.SetFields(rows)
		if err != nil {
			return nil, err
		}
		ans = append(ans, t)
	}
	if len(ans) == 0 {
		return nil, ErrNoRows
	}
	return ans, nil
}

func (s *Select[T]) Build() (*Query, error) {
	var t T
	m, err := s.core.r.Get(&t)
	if err != nil {
		return nil, err
	}
	s.model = m

	s.sb.WriteString("SELECT ")
	if err := s.buildColumns(); err != nil {
		return nil, err
	}
	s.sb.WriteString(" FROM ")
	if err = s.buildTable(s.table); err != nil {
		return nil, err
	}

	if len(s.where) > 0 {
		s.sb.WriteString(" WHERE ")
		pre := s.where[0]
		for i := 1; i < len(s.where); i++ {
			pre = pre.And(s.where[i])
		}
		err = s.BuildExpression(pre)
		if err != nil {
			return nil, err
		}
	}
	s.sb.WriteByte(';')

	return &Query{
		Sql:  s.sb.String(),
		Args: s.args,
	}, nil
}
func (s *Select[T]) From(tbl TableReference) *Select[T] {
	s.table = tbl
	return s
}
func (s *Select[T]) BuildExpression(p Expression) error {
	switch expr := p.(type) {
	case nil:
		return nil
	case *Column:
		err := s.buildColumn(expr, false)
		if err != nil {
			return err
		}
	case Value:
		s.sb.WriteByte('?')
		s.args = append(s.args, expr.val)
	case *Predicate:
		_, ok := expr.left.(*Predicate)
		if ok {
			s.sb.WriteString(" (")
		}
		err := s.BuildExpression(expr.left)
		if err != nil {
			return err
		}
		if ok {
			s.sb.WriteString(") ")
		}
		s.sb.WriteByte(' ')
		s.sb.WriteString(string(expr.op))
		s.sb.WriteByte(' ')
		_, ok = expr.right.(*Predicate)
		if ok {
			s.sb.WriteString(" (")
		}
		err = s.BuildExpression(expr.right)
		if err != nil {
			return err
		}
		if ok {
			s.sb.WriteString(") ")
		}
	case *SubQuery:
		s.sb.WriteByte('(')
		q, err := expr.Build()
		if err != nil {
			return err
		}
		s.sb.WriteString(q.Sql[:len(q.Sql)-1])

		s.sb.WriteByte(')')
	case SubQueryExpr:
		s.sb.WriteString(expr.pred)
		s.sb.WriteByte(' ')
		s.BuildExpression(expr.s)

	}
	return nil
}

func (s *Select[T]) Where(ps []*Predicate) *Select[T] {
	s.where = ps
	return s
}

func (s *Select[T]) Select(ss ...Selectable) *Select[T] {
	s.cols = ss
	return s
}

func NewSelect[T any](sess Session) *Select[T] {
	c := sess.getCore()
	return &Select[T]{
		typ:  "SELECT",
		core: c,
		sess: sess,
		builder: builder{
			core:    c,
			dialect: c.dialect,
		},
	}
}
func (s *Select[T]) AsSub(alias string) *SubQuery {
	var t T
	tbl := TableOf(&t)
	return &SubQuery{
		tbl:   tbl,
		sub:   s,
		alias: alias,
		cols:  s.cols,
	}
}
func (s *Select[T]) buildTable(t TableReference) error {
	switch tbl := t.(type) {
	case nil:
		s.builder.quota(s.model.TableName)
	case Table:
		m, err := s.r.Get(tbl.entity)
		if err != nil {
			return err
		}
		if tbl.alias != "" {
			s.quota(m.TableName)
			s.buildAs(tbl.alias)
		} else {
			s.builder.quota(m.TableName)
		}
	case Join:
		s.sb.WriteByte('(')
		err := s.buildTable(tbl.left)
		if err != nil {
			return err
		}
		s.sb.WriteByte(' ')
		s.sb.WriteString(tbl.typ)
		s.sb.WriteByte(' ')
		err = s.buildTable(tbl.right)
		if err != nil {
			return err
		}
		if len(tbl.on) > 0 {
			s.sb.WriteString(" ON ")
			pre := tbl.on[0]
			for i := 1; i < len(tbl.on); i++ {
				pre = pre.And(tbl.on[i])
			}
			err = s.BuildExpression(pre)
			if err != nil {
				return err
			}
		}
		if len(tbl.using) > 0 {
			s.sb.WriteString(" USING (")
			for i, col := range tbl.using {
				if i > 0 {
					s.sb.WriteByte(',')
				}
				err := s.buildColumn(&Column{name: col}, false)
				if err != nil {
					return err
				}
			}
			s.sb.WriteByte(')')
		}
		s.sb.WriteByte(')')
	case *SubQuery:
		s.sb.WriteByte('(')
		q, err := tbl.sub.Build()
		if err != nil {
			return err
		}
		s.sb.WriteString(q.Sql[:len(q.Sql)-1])
		s.args = append(s.args, q.Args...)
		s.sb.WriteByte(')')
		s.buildAs(tbl.alias)

	default:
		return errors.New("unknown table type")
	}
	return nil
}

func (s *Select[T]) buildColumns() error {
	if len(s.cols) == 0 {
		s.sb.WriteByte('*')
		return nil
	}
	for i, c := range s.cols {
		if i > 0 {
			s.sb.WriteByte(',')
		}
		switch val := c.(type) {
		case *Column:
			if err := s.buildColumn(val, true); err != nil {
				return err
			}
		case Aggregate:
			if err := s.buildAggregate(val, true); err != nil {
				return err
			}
		case RowExpr:
			s.sb.WriteString(val.expr)
			if len(val.args) > 0 {
				s.args = append(s.args, val.args...)
			}
		default:
			return errs.NewErrUnsupportedSelectable(val)

		}
	}
	return nil
}

func (s *Select[T]) buildColumn(col *Column, useAlias bool) error {
	err := s.builder.buildColumn(col.table, col.name)
	if err != nil {
		return err
	}
	if useAlias {
		s.buildAs(col.alias)
	}
	return nil
}
