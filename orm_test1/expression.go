package orm_test1

type RowExpr struct {
	expr string
	args []any
}

func (r RowExpr) selectable() {
}

func Raw(expr string, args ...any) RowExpr {
	return RowExpr{
		expr: expr,
		args: args,
	}
}
