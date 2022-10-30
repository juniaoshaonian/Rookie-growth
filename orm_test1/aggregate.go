package orm_test1

type Selectable interface {
	selectable()
	FieldName() string
}

type Aggregate struct {
	table TableReference
	fn    string
	name  string
	alias string
}

func (a Aggregate) FieldName() string {
	return a.name
}

func (a Aggregate) selectable() {
}

func (a Aggregate) Expr() {}

func (a Aggregate) Eq(v any) *Predicate {
	return &Predicate{
		left:  a,
		op:    Eq,
		right: Value{val: v},
	}
}

func (a Aggregate) Lt(v any) *Predicate {
	return &Predicate{
		left:  a,
		op:    Lt,
		right: Value{val: v},
	}
}

func (a Aggregate) Gt(v any) *Predicate {
	return &Predicate{
		left:  a,
		op:    Gt,
		right: Value{val: v},
	}
}

func (a Aggregate) Alias(alias string) Aggregate {

	return Aggregate{
		fn:    a.fn,
		name:  a.name,
		alias: alias,
	}
}

func Max(name string) Aggregate {
	return Aggregate{
		fn:   "MAX",
		name: name,
	}
}

func Min(name string) Aggregate {
	return Aggregate{
		fn:   "MIN",
		name: name,
	}
}

func Sum(name string) Aggregate {
	return Aggregate{
		fn:   "SUM",
		name: name,
	}
}

func Count(name string) Aggregate {
	return Aggregate{
		fn:   "COUNT",
		name: name,
	}
}

func Avg(name string) Aggregate {
	return Aggregate{
		fn:   "AVG",
		name: name,
	}
}
