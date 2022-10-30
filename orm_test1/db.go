package orm_test1

import (
	"context"
	"database/sql"
	"go.uber.org/multierr"
	"orm_test1/internal/valuer"
	"orm_test1/model"
)

type DB struct {
	db   *sql.DB
	core core
}

func (D *DB) getCore() core {
	return D.core
}

func (D *DB) queryContext(ctx context.Context, sql string, args []any) (*sql.Rows, error) {
	return D.db.QueryContext(ctx, sql, args...)
}

func (D DB) execContext(ctx context.Context, sql string, args []any) (sql.Result, error) {
	return D.db.ExecContext(ctx, sql, args)
}

type DBOption func(db *DB)

func Open(driver string, dsn string, opts ...DBOption) (*DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	return OpenDB(db, opts...)
}

func OpenDB(db *sql.DB, opts ...DBOption) (*DB, error) {
	d := &DB{
		db: db,
		core: core{
			r:       model.NewRegistry(),
			creator: valuer.NewReflectValuer,
			dialect: &mysqlDialect{},
		},
	}
	for _, opt := range opts {
		opt(d)
	}
	return d, nil
}

func DBWithCreator(creator valuer.ValuerCreater) DBOption {
	return func(db *DB) {
		db.core.creator = creator
	}
}

func (db *DB) Begin() (*Tx, error) {
	tx, err := db.db.Begin()
	if err != nil {
		return nil, err
	}
	return &Tx{
		tx:   tx,
		core: db.getCore(),
	}, nil
}

func (db *DB) DoTx(ctx context.Context, opts *sql.TxOptions, task func(ctx context.Context, tx *Tx) error) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	var paniced bool = true
	defer func() {
		if paniced || err != nil {
			err = multierr.Combine(err, tx.Rollback())
		} else {
			err = multierr.Combine(err, tx.Commit())
		}
	}()
	err = task(ctx, tx)
	return err
}

func DBWithMiddleware(ms ...Middleware) DBOption {
	return func(db *DB) {
		db.core.ms = ms
	}
}
