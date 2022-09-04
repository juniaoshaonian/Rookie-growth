package visitor

import (
	"ast/annotation"
	"errors"
	"go/ast"
)

// 获得interface的相关信息
type SingleFileVisitor struct {
	file *FileVisitor
}

type FileVisitor struct {
	Pkg   string
	anno  annotation.Annotations[*ast.File]
	types []*TypeVisitor
}

type TypeVisitor struct {
	Anno annotation.Annotations[*ast.TypeSpec]
	Name string
	T    Type
	I    *Interface
}
type Interface struct {
	Ann     annotation.Annotations[*ast.TypeSpec]
	Methods []*Method
}

type Method struct {
	Ann          annotation.Annotations[*ast.Field]
	Name         string
	Path         string
	ReqTypeName  []string
	RespTypeName []string
}

func (t *SingleFileVisitor) Visit(node ast.Node) (w ast.Visitor) {
	file, ok := node.(*ast.File)
	if ok {
		t.file = &FileVisitor{
			Pkg:  file.Name.Name,
			anno: annotation.NewAnnotations(file, file.Doc),
		}
		return t.file
	}
	return t
}

func (t *SingleFileVisitor) Get() (*FileVisitor, error) {
	if t.file == nil {
		return nil, errors.New("fileVisitor is nil")
	}
	return t.file, nil
}

func (f *FileVisitor) Visit(node ast.Node) (w ast.Visitor) {
	Type, ok := node.(*ast.TypeSpec)
	if ok {
		typvisitor := &TypeVisitor{Anno: annotation.NewAnnotations(Type, Type.Doc), Name: Type.Name.Name}
		f.types = append(f.types, typvisitor)
		return typvisitor
	}
	return f
}

func (f *FileVisitor) Get() ([]*TypeVisitor, error) {
	if len(f.types) == 0 {
		return nil, errors.New("interface is nil")
	}
	return f.types, nil
}

func (t *TypeVisitor) Visit(node ast.Node) (w ast.Visitor) {
	i, ok := node.(*ast.InterfaceType)
	if ok {
		T := Type{}
		T.Ann = t.Anno
		methods := i.Methods.List
		for i := 0; i < len(methods); i++ {
			f := Field{
				annotation.NewAnnotations(methods[i], methods[i].Doc),
			}
			T.Fields = append(T.Fields, f)
		}
		t.T = T
		t.I = &Interface{}
	}
	return t
}

func (t *TypeVisitor) Get() (*Interface, error) {
	if t.I != nil {

		return t.I, nil
	}
	return nil, errors.New("interface is nil")
}

type Type struct {
	Ann    annotation.Annotations[*ast.TypeSpec]
	Fields []Field
}

type Field struct {
	annotation.Annotations[*ast.Field]
}
