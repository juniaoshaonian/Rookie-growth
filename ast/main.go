package main

import (
	"ast/http"
	"ast/visitor"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

// 实际上 main 函数这里要考虑接收参数
// src 源目标
// dst 目标目录
// type src 里面可能有很多类型，那么用户可能需要指定具体的类型
// 这里我们简化操作，只读取当前目录下的数据，并且扫描下面的所有源文件，然后生成代码
// 在当前目录下运行 go install 就将 main 安装成功了，
// 可以在命令行中运行 gen
// 在 testdata 里面运行 gen，则会生成能够通过所有测试的代码

var (
	ErrIgnoreType error = errors.New("ignore this type")
)

func main() {

	srcFiles, err := scanFiles("./testdata")
	if err != nil {
		fmt.Println(err)
	}
	fset := token.NewFileSet()
	fv, err := parser.ParseFile(fset, srcFiles[1], nil, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
	}
	s := &visitor.SingleFileVisitor{}
	ast.Walk(s, fv)
	//parseFiles(srcFiles)
	fmt.Println("success")
}

func gen(src string) error {
	// 第一步找出符合条件的文件
	srcFiles, err := scanFiles(src)
	if err != nil {
		return err
	}
	// 第二步，AST 解析源代码文件，拿到 service definition 定义
	defs, err := parseFiles(srcFiles)
	if err != nil {
		return err
	}
	// 生成代码
	return genFiles(src, defs)
}

// 根据 defs 来生成代码
// src 是源代码所在目录，在测试里面它是 ./testdata
func genFiles(src string, defs []http.ServiceDefinition) error {
	for _, s := range defs {
		filename := underscoreName(s.Name) + "_gen" + ".go"
		fmt.Println(src + "/" + filename)
		f, err := os.Create(src + "/" + filename)
		if err != nil {
			return err
		}
		err = http.Gen(f, s)
		if err != nil {
			return err
		}
	}
	return nil
}

func parseFiles(srcFiles []string) ([]http.ServiceDefinition, error) {
	ans := make([]http.ServiceDefinition, 0, len(srcFiles))
	for _, src := range srcFiles {
		fset := token.NewFileSet()
		f, _ := parser.ParseFile(fset, src, nil, parser.ParseComments)
		tv := &visitor.SingleFileVisitor{}
		ast.Walk(tv, f)
		filev, err := tv.Get()
		if err != nil {
			return nil, err
		}
		pkg := filev.Pkg
		typv, err := filev.Get()
		if err != nil {
			return nil, err
		}
		for _, t := range typv {
			if t.I == nil {
				continue
			}
			var HttpClient bool
			for _, a := range t.Anno.Ans {
				if a.Key == "HttpClient" {
					HttpClient = true
				}
			}
			if !HttpClient {
				continue
			}
			sd, err := parseServiceDefinition(pkg, t.T)
			if err != nil {
				if err == ErrIgnoreType {
					continue
				}
			}
			ans = append(ans, sd)

		}
	}
	return ans, nil
}

// 返回符合条件的 Go 源代码文件，也就是你要用 AST 来分析这些文件的代码
func scanFiles(src string) ([]string, error) {
	files, err := os.ReadDir(src)
	if err != nil {
		return []string{}, err
	}
	filestxt := []string{}
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".go") {
			continue
		}
		filetxt, err := filepath.Abs(src + "/" + file.Name())
		if err != nil {
			return []string{}, err
		}
		filestxt = append(filestxt, filetxt)
	}
	return filestxt, nil
}

// underscoreName 驼峰转字符串命名，在决定生成的文件名的时候需要这个方法
// 可以用正则表达式，然而我写不出来，我是正则渣
func underscoreName(name string) string {
	var buf []byte
	for i, v := range name {
		if unicode.IsUpper(v) {
			if i != 0 {
				buf = append(buf, '_')
			}
			buf = append(buf, byte(unicode.ToLower(v)))
		} else {
			buf = append(buf, byte(v))
		}

	}
	return string(buf)
}

func parseServiceDefinition(pkg string, typ visitor.Type) (http.ServiceDefinition, error) {

	sd := http.ServiceDefinition{}
	sd.Package = pkg
	for _, a := range typ.Ann.Ans {
		if a.Key == "ServiceName" {
			sd.Name = a.Value
		}

	}
	if sd.Name == "" {
		sd.Name = typ.Ann.Node.Name.Name
	}

	InT := visitor.Interface{Ann: typ.Ann}
	for _, f := range typ.Fields {
		m := &visitor.Method{}
		m.Name = f.Node.Names[0].Name
		params := f.Node.Type.(*ast.FuncType).Params.List
		for _, a := range f.Ans {
			if a.Key == "Path" {

				m.Path = a.Value
			}
		}
		if m.Path == "" {
			m.Path = fmt.Sprintf("/%s", m.Name)
		}

		for _, param := range params {
			switch a := param.Type.(type) {
			case *ast.SelectorExpr:
				packagename := a.X.(*ast.Ident).Name
				funcname := a.Sel.Name
				m.ReqTypeName = append(m.ReqTypeName, packagename+"."+funcname)
			case *ast.StarExpr:
				funcname := a.X.(*ast.Ident).Name
				m.ReqTypeName = append(m.ReqTypeName, funcname)
			}
		}
		if f.Node.Type.(*ast.FuncType).Results == nil {
			InT.Methods = append(InT.Methods, m)
			continue
		}
		results := f.Node.Type.(*ast.FuncType).Results.List
		for _, res := range results {
			switch a := res.Type.(type) {
			case *ast.SelectorExpr:
				packagename := a.X.(*ast.Ident).Name
				funcname := a.Sel.Name
				m.RespTypeName = append(m.RespTypeName, packagename+"."+funcname)
			case *ast.StarExpr:
				funcname := a.X.(*ast.Ident).Name
				m.RespTypeName = append(m.RespTypeName, funcname)
			case *ast.Ident:
				funcname := a.Name
				m.RespTypeName = append(m.RespTypeName, funcname)
			}
		}
		InT.Methods = append(InT.Methods, m)
	}

	for _, m := range InT.Methods {

		sm := http.ServiceMethod{}
		sm.Name = m.Name
		sm.Path = m.Path
		var Contextflag bool
		for _, req := range m.ReqTypeName {
			if req == "context.Context" {
				Contextflag = true
				continue
			}
			sm.ReqTypeName = req
		}
		if Contextflag == false {
			return http.ServiceDefinition{}, errors.New("gen: 方法必须接收两个参数，其中第一个参数是 context.Context，第二个参数请求")
		}
		var errorFlag bool
		for _, resp := range m.RespTypeName {
			if resp == "error" {
				errorFlag = true
				continue
			}
			sm.RespTypeName = resp
		}
		if errorFlag == false {
			return http.ServiceDefinition{}, errors.New("gen: 方法必须返回两个参数，其中第一个返回值是响应，第二个返回值是error")
		}
		sd.Methods = append(sd.Methods, sm)

	}

	return sd, nil

}
