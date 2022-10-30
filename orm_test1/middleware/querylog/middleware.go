package querylog

import (
	"context"
	"orm_test1"
)

type MiddlewareBuilder struct {
	logFunc func(sql string, args []any)
}

func (m *MiddlewareBuilder) LogFunc(logfunc func(sql string, args []any)) *MiddlewareBuilder {
	m.logFunc = logfunc
	return m
}
func (m *MiddlewareBuilder) Build() orm_test1.Middleware {
	return func(next orm_test1.Handler) orm_test1.Handler {
		return func(ctx context.Context, qc *orm_test1.QueryContext) *orm_test1.QueryResult {
			q, err := qc.Builder.Build()
			if err != nil {
				return &orm_test1.QueryResult{
					Err: err,
				}
			}
			m.logFunc(q.Sql, q.Args)
			return next(ctx, qc)
		}
	}
}
