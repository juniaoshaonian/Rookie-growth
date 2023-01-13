package model

type Model struct {
	Fields   []*Field
	FieldMap map[string]*Field
}

type Field struct {
	ColName   string
	FieldName string
}
