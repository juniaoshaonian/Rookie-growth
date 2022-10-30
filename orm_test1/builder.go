package orm_test1

import (
	"errors"
	"orm_test1/internal/errs"
	"orm_test1/model"
	"strings"
)

type builder struct {
	core
	model   *model.Model
	sb      strings.Builder
	args    []any
	dialect Dialect
}

func (b *builder) quota(name string) {
	b.sb.WriteByte(b.dialect.quoter())
	b.sb.WriteString(name)
	b.sb.WriteByte(b.dialect.quoter())
}

func (b *builder) buildColumn(table TableReference, fd string) error {

	var alias string
	if table != nil {
		alias = table.tableAlias()
	}
	if alias != "" {
		b.quota(alias)
		b.sb.WriteByte('.')
	}
	colName, err := b.colName(table, fd)
	if err != nil {
		return err
	}
	b.quota(colName)
	return nil
}

func (b *builder) colName(table TableReference, fd string) (string, error) {
	switch tbl := table.(type) {
	case nil:
		f, ok := b.model.FieldMap[fd]
		if !ok {
			return "", errs.NewErrField(fd)
		}
		return f.Colname, nil
	case Table:
		m, err := b.r.Get(tbl.entity)
		if err != nil {
			return "", err
		}
		fdmeta, ok := m.FieldMap[fd]
		if !ok {
			return "", errs.NewErrField(fd)
		}
		return fdmeta.Colname, nil
	case Join:
		colName, err := b.colName(tbl.left, fd)
		if err == nil {
			return colName, nil
		}
		return b.colName(tbl.right, fd)
	case *SubQuery:
		if len(tbl.cols) > 0 {
			for _, col := range tbl.cols {
				if col.FieldName() != fd {
					return "", errs.NewErrField(fd)
				}
			}
		}
		return b.colName(tbl.tbl, fd)
	default:
		return "", errors.New("错误的表类型")
	}
}
func (b *builder) buildAggregate(a Aggregate, useAlias bool) error {
	b.sb.WriteString(a.fn)
	b.sb.WriteByte('(')
	err := b.buildColumn(a.table, a.name)
	if err != nil {
		return err
	}
	if useAlias {
		b.buildAs(a.alias)
	}
	return nil
}
func (b *builder) buildAs(alias string) {
	if alias != "" {
		b.sb.WriteString(" AS ")
		b.quota(alias)
	}
}
