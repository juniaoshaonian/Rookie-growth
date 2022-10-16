package orm

import (
	"orm/internal/errs"
)

type Dialect interface {
	quoter() byte
	buildDuplicateKey(b *builder, odk *OnDuplicateKey) error
}

type standardSQL struct {
}

type mysqlDialct struct {
	standardSQL
}

func (m mysqlDialct) quoter() byte {
	//TODO implement me
	panic("implement me")
}

func (d *mysqlDialct) buildDuplicateKey(b *builder, odk *OnDuplicateKey) error {

	b.sb.WriteString(" ON DUPLICATE KEY UPDATE ")
	for idx, assign := range odk.assigns {
		if idx > 0 {
			b.sb.WriteByte(',')
		}
		switch expr := assign.(type) {
		case Assignment:
			fd, ok := b.model.Fieldmap[expr.column]
			if !ok {
				return errs.NewErrUnknownField(expr.column)
			}
			b.sb.WriteByte('`')
			b.sb.WriteString(fd.Colname)
			b.sb.WriteByte('`')
			b.sb.WriteString("=?")
			b.args = append(b.args, expr.val)
		case Column:
			fd, ok := b.model.Fieldmap[expr.name]
			if !ok {
				return errs.NewErrUnknownField(expr.name)
			}
			b.sb.WriteByte('`')
			b.sb.WriteString(fd.Colname)
			b.sb.WriteByte('`')
			b.sb.WriteString("=VALUES(")
			b.sb.WriteByte('`')
			b.sb.WriteString(fd.Colname)
			b.sb.WriteByte('`')
			b.sb.WriteByte(')')

		}
	}
	return nil
}

type sqliteDialect struct {
	standardSQL
}
