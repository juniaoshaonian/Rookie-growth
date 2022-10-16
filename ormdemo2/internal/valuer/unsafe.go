package valuer

import (
	"database/sql"
)

type unsafeValuer struct {
}

func (u unsafeValuer) SetColumns(rows *sql.Rows) error {
	//TODO implement me
	panic("implement me")
}
