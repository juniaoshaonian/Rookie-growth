package zorm

import "context"

type Select[T any] struct {
	builder
	table string
}

func (s *Select[T]) Get(ctx context.Context) (T, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Select[T]) Build() (Query, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Select) From()
