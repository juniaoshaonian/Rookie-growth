package orm_test1

import "context"

type QueryBuilder interface {
	Build() (*Query, error)
}
type Querier[T any] interface {
	Get(ctx context.Context) (*T, error)
	GetMulti(ctx context.Context) ([]*T, error)
}

type Executor interface {
	Exec(ctx context.Context) Result
}

type Query struct {
	Sql  string
	Args []any
}
