package orm

import (
	"database/sql"
	"orm/internal/model"
	"orm/internal/valuer"
)

type DB struct {
	db         *sql.DB
	r          model.Registery
	valCreator valuer.Creater
	dialect    Dialect
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
	new_DB := &DB{
		r:          model.NewRegitery(),
		db:         db,
		valCreator: valuer.NewReflectValue,
		dialect:    &mysqlDialct{},
	}
	for _, opt := range opts {
		opt(new_DB)
	}
	return new_DB, nil
}

func DBWithRegistry(r model.Registery) DBOption {
	return func(db *DB) {
		db.r = r
	}
}

func DBWithDialect(dialect Dialect) DBOption {
	return func(db *DB) {
		db.dialect = dialect
	}
}
