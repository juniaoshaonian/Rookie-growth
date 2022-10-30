package valuer

import (
	"database/sql"
	"orm_test1/model"
)

type Valuer interface {
	SetFields(rows *sql.Rows) error
	Field(name string) (any, error)
}

type ValuerCreater func(m *model.Model, t any) Valuer
