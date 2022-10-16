package valuer

import (
	"database/sql"
	"orm/internal/model"
)

type Valuer interface {
	SetColumns(rows *sql.Rows) error
}

type Creater func(t any, model *model.Model) Valuer
