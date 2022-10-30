package orm_test1

import (
	"context"
	"database/sql"
	"orm_test1/internal/valuer"
	"orm_test1/model"
)

type Tx struct {
	tx *sql.Tx
	core
}

func (tx *Tx) getCore() core {
	return tx.core
}

func (tx *Tx) queryContext(ctx context.Context, sql string, args []any) (*sql.Rows, error) {
	return tx.tx.QueryContext(ctx, sql, args)
}

func (tx *Tx) execContext(ctx context.Context, sql string, args []any) (sql.Result, error) {
	return tx.tx.ExecContext(ctx, sql, args)
}

func (tx *Tx) Commit() error {
	return tx.tx.Commit()
}

func (tx *Tx) Rollback() error {
	return tx.tx.Rollback()
}

type Session interface {
	getCore() core
	queryContext(ctx context.Context, sql string, args []any) (*sql.Rows, error)
	execContext(ctx context.Context, sql string, args []any) (sql.Result, error)
}

type core struct {
	r       model.Registry
	creator valuer.ValuerCreater
	dialect Dialect
	ms      []Middleware
}
