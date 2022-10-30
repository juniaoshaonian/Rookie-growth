package orm_test1

type Expression interface {
	Expr()
}

type Predicate struct {
	left  Expression
	op    Op
	right Expression
}

func (p *Predicate) Expr() {

}

func (p *Predicate) And(p1 *Predicate) *Predicate {
	return &Predicate{
		left:  p,
		op:    And,
		right: p1,
	}
}

func (p *Predicate) Or(p1 *Predicate) *Predicate {
	return &Predicate{
		left:  p,
		op:    OR,
		right: p1,
	}
}

func (p *Predicate) Not() *Predicate {
	return &Predicate{
		op:    NOT,
		right: p,
	}

}
