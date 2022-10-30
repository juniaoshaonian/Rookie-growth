package orm_test1

type Assignable interface {
	assign()
}

type Assignment struct {
	column string
	val    Value
}

func (a Assignment) assign() {}

func Assign(col string, val any) Assignment {
	return Assignment{
		column: col,
		val: Value{
			val: val,
		},
	}

}
