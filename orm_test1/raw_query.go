package orm_test1

import "context"

type RawQuerier[T any] struct {
	typ string
	core
	sess Session
	sql  string
	args []any
}

func (r *RawQuerier[T]) Build() (*Query, error) {
	return &Query{
		Sql:  r.sql,
		Args: r.args,
	}, nil
}

func RawQuery[T any](sess Session, sql string, args ...any) *RawQuerier[T] {
	return &RawQuerier[T]{
		sql:  sql,
		args: args,
		typ:  "RAW",
		sess: sess,
		core: sess.getCore(),
	}
}

func (r *RawQuerier[T]) Get(ctx context.Context) (*T, error) {
	model, err := r.r.Get(new(T))
	if err != nil {
		return nil, err
	}
	res := get[T](ctx, r.core, r.sess, &QueryContext{
		Type:      r.typ,
		Builder:   r,
		Model:     model,
		TableName: model.TableName,
	})
	if res.Res != nil {
		return res.Res.(*T), nil
	}
	return nil, res.Err
}

func get[T any](ctx context.Context, c core, sess Session, qc *QueryContext) *QueryResult {
	root := func(ctx context.Context, qc *QueryContext) *QueryResult {
		q, err := qc.Builder.Build()
		if err != nil {
			return &QueryResult{
				Err: err,
			}
		}
		rows, err := sess.queryContext(ctx, q.Sql, q.Args)
		if err != nil {
			return &QueryResult{
				Err: err,
			}
		}
		t := new(T)
		m, err := c.r.Get(t)
		if err != nil {
			return &QueryResult{
				Err: err,
			}
		}
		valuerGet := c.creator(m, t)
		if !rows.Next() {
			return &QueryResult{
				Err: ErrNoRows,
			}
		}
		err = valuerGet.SetFields(rows)
		if err != nil {
			return &QueryResult{
				Err: err,
			}
		}
		return &QueryResult{
			Res: t,
		}
	}
	for i := len(c.ms) - 1; i >= 0; i-- {
		root = c.ms[i](root)
	}
	res := root(ctx, qc)
	if res.Res != nil {
		return &QueryResult{
			Res: res,
		}
	}
	return &QueryResult{
		Err: res.Err,
	}
}
