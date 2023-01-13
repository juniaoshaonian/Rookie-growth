package zorm

import "context"

type Builder interface {
	Build() (Query, error)
}
type Querier[T any] interface {
	Get(ctx context.Context) (T, error)
}

type Query struct {
	Sql  string
	Args []any
}
