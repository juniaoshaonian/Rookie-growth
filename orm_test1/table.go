package orm_test1

type TableReference interface {
	tableAlias() string
}

// 普通表
type Table struct {
	entity any
	alias  string
}

func TableOf(entity any) Table {
	return Table{
		entity: entity,
	}
}

func (t Table) As(alias string) Table {
	return Table{
		entity: t.entity,
		alias:  alias,
	}
}
func (t Table) tableAlias() string {
	return t.alias
}

func (t Table) Join(right TableReference) *JoinBuilder {
	return &JoinBuilder{
		left:  t,
		typ:   "JOIN",
		right: right,
	}
}
func (t Table) LeftJoin(right TableReference) *JoinBuilder {
	return &JoinBuilder{
		left:  t,
		typ:   "LEFT JOIN",
		right: right,
	}
}
func (t Table) RightJoin(right TableReference) *JoinBuilder {
	return &JoinBuilder{
		left:  t,
		typ:   "RIGHT JOIN",
		right: right,
	}
}
func (t Table) C(name string) *Column {
	return &Column{
		table: t,
		name:  name,
	}
}

// join查询
type Join struct {
	left  TableReference
	typ   string
	right TableReference
	on    []*Predicate
	using []string
}

func (j Join) Join(target TableReference) *JoinBuilder {
	return &JoinBuilder{
		left:  j,
		right: target,
		typ:   "JOIN",
	}
}

func (j Join) LeftJoin(target TableReference) *JoinBuilder {
	return &JoinBuilder{
		left:  j,
		right: target,
		typ:   "LEFT JOIN",
	}
}

func (j Join) RightJoin(target TableReference) *JoinBuilder {
	return &JoinBuilder{
		left:  j,
		right: target,
		typ:   "RIGHT JOIN",
	}
}

func (j Join) tableAlias() string {
	return ""
}

type JoinBuilder struct {
	left  TableReference
	typ   string
	right TableReference
}

func (jb *JoinBuilder) On(ps ...*Predicate) Join {
	return Join{
		left:  jb.left,
		typ:   jb.typ,
		right: jb.right,
		on:    ps,
	}
}

func (jb *JoinBuilder) Using(cols ...string) Join {
	return Join{
		left:  jb.left,
		typ:   jb.typ,
		right: jb.right,
		using: cols,
	}
}

// 子查询
type SubQuery struct {
	alias string
	sub   QueryBuilder
	tbl   Table
	q     *Query
	cols  []Selectable
}

func (s *SubQuery) Build() (q *Query, err error) {
	if s.q == nil {
		q, err = s.sub.Build()
		s.q = q
	}
	return s.q, err
}
func (s *SubQuery) Expr() {
}

func (s *SubQuery) tableAlias() string {
	return s.alias
}

func (s *SubQuery) C(name string) *Column {
	return &Column{
		table: s,
		name:  name,
	}
}

func (s *SubQuery) Join(table TableReference) *JoinBuilder {
	return &JoinBuilder{
		left:  s,
		typ:   "JOIN",
		right: table,
	}
}
func (s *SubQuery) RightJoin(table TableReference) *JoinBuilder {
	return &JoinBuilder{
		left:  s,
		typ:   "RIGHT JOIN",
		right: table,
	}
}
func (s *SubQuery) LeftJoin(table TableReference) *JoinBuilder {
	return &JoinBuilder{
		left:  s,
		typ:   "LEFT JOIN",
		right: table,
	}
}
