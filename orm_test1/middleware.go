package orm_test1

import (
	"context"
	"orm_test1/model"
)

type QueryContext struct {
	Type      string
	Builder   QueryBuilder
	Model     *model.Model
	TableName string
}
type QueryResult struct {
	Res any
	Err error
}

type Handler func(ctx context.Context, qc *QueryContext) *QueryResult

type Middleware func(next Handler) Handler
