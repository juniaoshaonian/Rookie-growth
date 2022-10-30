package orm_test1

type Op string

const (
	Eq    Op = "="
	Lt    Op = "<"
	Gt    Op = ">"
	And   Op = "AND"
	OR    Op = "OR"
	NOT   Op = "NOT"
	IN    Op = "IN"
	EXIST Op = "EXIST"
)

type Column struct {
	table TableReference
	name  string
	alias string
}

func (c *Column) assign() {}

func (c *Column) selectable() {}

func C(name string) *Column {
	return &Column{
		name: name,
	}
}

type Value struct {
	val any
}

func (v Value) Expr() {}

func (c *Column) Expr() {}

func (c *Column) EQ(val any) *Predicate {
	return &Predicate{
		left:  c,
		op:    Eq,
		right: exprOf(val),
	}
}
func (c *Column) LT(val any) *Predicate {
	return &Predicate{
		left:  c,
		op:    Lt,
		right: exprOf(val),
	}
}

func (c *Column) GT(val any) *Predicate {
	return &Predicate{
		left:  c,
		op:    Gt,
		right: exprOf(val),
	}
}

func (c *Column) Alias(alias string) *Column {
	c.alias = alias
	return c
}

func exprOf(e any) Expression {
	switch exp := e.(type) {
	case Expression:
		return exp
	default:
		return Value{
			val: exp,
		}
	}
}
func (c *Column) In(sub *SubQuery) *Predicate {
	return &Predicate{
		left:  c,
		op:    IN,
		right: sub,
	}
}
func Exist(sub *SubQuery) *Predicate {
	return &Predicate{
		op:    EXIST,
		right: sub,
	}
}

type SubQueryExpr struct {
	s    *SubQuery
	pred string
}

func (s SubQueryExpr) Expr() {

}
func ALL(sub *SubQuery) SubQueryExpr {
	return SubQueryExpr{
		s:    sub,
		pred: "ALL",
	}
}

func SOME(sub *SubQuery) SubQueryExpr {
	return SubQueryExpr{
		s:    sub,
		pred: "SOME",
	}
}
func ANY(sub *SubQuery) SubQueryExpr {
	return SubQueryExpr{
		s:    sub,
		pred: "ANY",
	}
}
