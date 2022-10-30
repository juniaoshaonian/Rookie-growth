package orm_test1

import (
	"orm_test1/internal/errs"
)

type Dialect interface {
	quoter() byte
	buildDuplicateKey(sb *builder, odk *OnDuplicateKey) error
}

type standardSQL struct {
}

type mysqlDialect struct {
	standardSQL
}

func (m *mysqlDialect) quoter() byte {
	return '`'
}

func (d *mysqlDialect) buildDuplicateKey(b *builder, odk *OnDuplicateKey) error {
	b.sb.WriteString(" ON DUPLICATE KEY UPDATE ")
	for idx, assign := range odk.assigns {
		if idx > 0 {
			b.sb.WriteByte(',')
		}
		switch expr := assign.(type) {
		case Assignment:
			fd, ok := b.model.FieldMap[expr.column]
			if !ok {
				return errs.NewErrField(expr.column)
			}
			b.quota(fd.Colname)
			b.sb.WriteByte('=')
			b.sb.WriteByte('?')
			b.args = append(b.args, expr.val.val)
		case *Column:
			fd, ok := b.model.FieldMap[expr.name]
			if !ok {
				return errs.NewErrField(expr.name)
			}
			b.quota(fd.Colname)
			b.sb.WriteString("=VALUES(")
			b.quota(fd.Colname)
			b.sb.WriteString(")")
		}
	}
	return nil
}

type sqliteDialect struct {
	standardSQL
}
