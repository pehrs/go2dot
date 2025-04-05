package golang

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"strings"
	"unicode"
)

type GolangStructField struct {
	Ast  *ast.Field
	Name string `json:"name"`
	Type string `json:"type"`
}

type GolangFuncArg struct {
	Name string
	Type string
}

type GolangStructFunc struct {
	Ast        *ast.FuncDecl
	Name       string
	ReturnType string
	Args       []GolangFuncArg
}

type GolangStruct struct {
	Ast    *ast.StructType
	Name   string              `json:"name"`
	Fields []GolangStructField `json:"fields"`
	Funcs  []GolangStructFunc  `json:"functions"`
}

func (str GolangStruct) String() string {
	j, _ := json.MarshalIndent(str, "", "  ")
	return string(j)
}

func (str GolangStruct) TypeId() string {
	return getTypeId(str.Ast)
}

func (str GolangStruct) DotDeps(structName map[(*GolangStruct)]string, pool map[string](*GolangStruct)) string {

	result := ""

	for _, field := range str.Fields {

		typeId := getTypeId(field.Ast.Type)

		relatedStruct, ok := pool[typeId]
		if ok {
			relatedName := structName[relatedStruct]
			fmt.Printf("\"%v\" -> \"%v\" [arrowhead=open style=dashed];\n", str.Name, relatedName)
		}
	}

	return result
}

//	 Interface1[
//		label = <{<b>«interface» I/O</b> | + property<br/>...<br/>|+ method<br/>...<br/>}>,
//		shape=record
//
// ];
func (str GolangStruct) Dot() string {

	result := fmt.Sprintf("%s[label = <{<b>«struct»<br/>%s</b><br align=\"left\"/>|",
		str.Name, str.Name)

	for _, field := range str.Fields {
		access := "+"
		if unicode.IsLower([]rune(field.Name)[0]) {
			access = "-"
		}
		result = fmt.Sprintf("%s%s %s %s<br align=\"left\"/>",
			result, access, field.Name, dotEscape(field.Type),
		)
	}
	result = fmt.Sprintf("%s|", result)
	for _, theFunc := range str.Funcs {
		access := "+"
		if unicode.IsLower([]rune(theFunc.Name)[0]) {
			access = "-"
		}
		paramDecl := "param"
		if theFunc.Ast.Type.Params != nil {
			paramDecl = fmt.Sprintf("(%s)", fields(*theFunc.Ast.Type.Params))
		}
		returnDecl := "res"
		// print return params
		if theFunc.Ast.Type.Results != nil {
			if len(theFunc.Ast.Type.Results.List) > 1 {
				returnDecl = fmt.Sprintf("(%s)", fields(*theFunc.Ast.Type.Results))
			} else {
				returnDecl = fmt.Sprintf("%s", fields(*theFunc.Ast.Type.Results))
			}
		}
		result = fmt.Sprintf("%s%s %s%s %s<br align=\"left\"/>",
			result, access, theFunc.Name, dotEscape(paramDecl), dotEscape(returnDecl),
		)
	}

	result = fmt.Sprintf("%s}>, shape=record];", result)

	return result
}

func NewGolangStruct(name string, str *ast.StructType) *GolangStruct {

	goStruct := new(GolangStruct)

	goStruct.Name = name
	goStruct.Ast = str

	for _, field := range str.Fields.List {

		for _, fieldName := range field.Names {
			goStruct.Fields = append(goStruct.Fields, GolangStructField{
				Ast:  field,
				Name: fieldName.String(),
				Type: expr(field.Type),
			})
		}

	}

	return goStruct
}

func getAstFile(filename string) (*ast.File, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	srcbuf, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	src := string(srcbuf)

	// file set
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, "lib.go", src, 0)
	if err != nil {
		return nil, err
	}
	return astFile, nil
}

func ExtractStructs(filename string) ([]*GolangStruct, error) {
	astFile, err := getAstFile(filename)
	if err != nil {
		return nil, err
	}

	var result = []*GolangStruct{}

	var structMap = map[string]*GolangStruct{}
	var structFuncMap = map[string][]GolangStructFunc{}

	ast.Inspect(astFile, func(n ast.Node) bool {

		switch fn := n.(type) {
		case *ast.TypeSpec:
			str, ok := fn.Type.(*ast.StructType)
			if ok {
				//  Interface1[
				// 	label = <{<b>«interface» I/O</b> | + property<br/>...<br/>|+ method<br/>...<br/>}>,
				// 	shape=record
				// ];

				// result = fmt.Sprintf("%s\n%s [label = <{<b>«struct»<br/>%s</b>| %s}>]",
				// 	result, fn.Name.String(), fn.Name.String(), structDotFields(*str.Fields))

				theStruct := NewGolangStruct(fn.Name.String(), str)

				result = append(result, theStruct)
				structMap[fn.Name.String()] = theStruct
			}
		case *ast.FuncDecl:

			// We are only looking for singler reciever functions
			if fn.Recv != nil && len(fn.Recv.List) == 1 {
				recv := fn.Recv.List[0]
				structName := getTypeId(recv.Type)

				returnType := ""
				if fn.Type.Results != nil {
					if len(fn.Type.Results.List) > 1 {
						returnType = fmt.Sprintf("(%s)", fields(*fn.Type.Results))
					} else {
						returnType = fmt.Sprintf("%s", fields(*fn.Type.Results))
					}
				}

				// theFunc := new(GolangStructFunc)
				// theFunc.Ast = fn
				// theFunc.Name = structName
				// theFunc.ReturnType = returnType
				// theFunc.Args = []GolangFuncArg{}

				theFunc := GolangStructFunc{
					Ast:        fn,
					Name:       fn.Name.String(),
					ReturnType: returnType,
					Args:       []GolangFuncArg{},
				}

				for _, param := range fn.Type.Params.List {
					for _, paramName := range param.Names {
						theFunc.Args = append(theFunc.Args, GolangFuncArg{
							Name: paramName.String(),
							Type: expr(param.Type),
						})
					}
				}
				list, ok := structFuncMap[structName]
				if ok {
					structFuncMap[structName] = append(list, theFunc)
				} else {
					structFuncMap[structName] = []GolangStructFunc{
						theFunc,
					}
				}

			}
			// // if a method, explore and print receiver
			// if fn.Recv != nil {
			// 	fmt.Printf("(%s)", fields(*fn.Recv))
			// }

			// // print actual function name
			// fmt.Printf("%v", fn.Name)

			// // print function parameters
			// if fn.Type.Params != nil {
			// 	fmt.Printf("(%s)", fields(*fn.Type.Params))
			// }

			// // print return params
			// if fn.Type.Results != nil {
			// 	fmt.Printf("(%s)", fields(*fn.Type.Results))
			// }

			fmt.Sprintf("%v", fn)

		}

		for structName, funcList := range structFuncMap {
			theStruct, ok := structMap[structName]
			if ok {
				theStruct.Funcs = funcList
			}
		}

		return true
	})

	return result, nil
}

func expr(e ast.Expr) (ret string) {
	switch x := e.(type) {
	case *ast.StarExpr:
		sel, ok := x.X.(*ast.SelectorExpr)
		if ok {
			return fmt.Sprintf("%s*%s.%s",
				ret,
				sel.X.(*ast.Ident).Name,
				sel.Sel.Name,
			)
		}
		// x.X.(*ast.SelectorExpr).X.(*ast.Ident).Name
		// x.X.(*ast.SelectorExpr).Sel.(*ast.Ident).Name
		return fmt.Sprintf("%s*%v", ret, x.X)
	case *ast.Ident:
		return fmt.Sprintf("%s%v", ret, x.Name)
	case *ast.ArrayType:
		if x.Len != nil {
			return "TODO: ARRAY"
		}
		res := expr(x.Elt)
		return fmt.Sprintf("%s[]%v", ret, res)
	case *ast.MapType:
		return fmt.Sprintf("map[%s]%s", expr(x.Key), expr(x.Value))
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", expr(x.X), expr(x.Sel))
	case *ast.InterfaceType:
		return "inteface{}"
	case *ast.StructType:
		return "struct"
	case *ast.FuncType:
		returnType := ""
		if x.Results != nil {
			returnType = fields(*x.Results)
		}
		params := ""
		if x.Params != nil {
			if len(x.Params.List) > 1 {
				params = fmt.Sprintf("(%s)", fields(*x.Params))
			} else {
				params = fmt.Sprintf("%s", fields(*x.Params))
			}
		}
		return fmt.Sprintf("func(%s) %s", params, returnType)

	case *ast.ParenExpr:
		return fmt.Sprintf("(%s)", expr(x.X))
	case *ast.Ellipsis:
		res := expr(x.Elt)
		return fmt.Sprintf("%s%v...", ret, res)
	default:
		fmt.Printf("\nUNKOWN: %#v\n", x)
	}
	return
}

func getTypeId(e ast.Expr) string {
	switch x := e.(type) {
	case *ast.StarExpr:
		sel, ok := x.X.(*ast.SelectorExpr)
		if ok {
			return fmt.Sprintf("%s.%s",
				sel.X.(*ast.Ident).Name,
				sel.Sel.Name,
			)
		}

		return fmt.Sprintf("%v", x.X)
	case *ast.Ident:
		return fmt.Sprintf("%v", x.Name)
	case *ast.ArrayType:
		if x.Len != nil {
			return "TODO: ARRAY"
		}
		res := getTypeId(x.Elt)
		return fmt.Sprintf("%v", res)
	case *ast.MapType:
		return fmt.Sprintf("%s", getTypeId(x.Value))
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", expr(x.X), expr(x.Sel))
	case *ast.InterfaceType:
		return "inteface{}"
	case *ast.StructType:
		return "struct"
	case *ast.FuncType:
		return "func"
	default:
		fmt.Printf("\nTODO UNKOWN: %#v\n", x)
	}
	return "?"
}

func dotEscape(in string) string {
	res := strings.Replace(in, "&", "&amp;", -1)
	res = strings.Replace(res, "{", "&#123;", -1)
	res = strings.Replace(res, "}", "&#125;", -1)

	return res
}

func fields(fl ast.FieldList) (ret string) {
	pcomma := ""
	for i, f := range fl.List {
		// get all the names if present
		var names string
		ncomma := ""
		for j, n := range f.Names {
			if j > 0 {
				ncomma = ", "
			}
			names = fmt.Sprintf("%s%s%s ", names, ncomma, n)
		}
		if i > 0 {
			pcomma = ", "
		}
		ret = fmt.Sprintf("%s%s%s%s", ret, pcomma, names, expr(f.Type))
	}
	return ret
}
