package samples

var PkgStringVar = "value"
var pkgIntVar = 42
var pkgFloatVar = 42.42

type ChildStruct struct {
	name string
}

type ParentStruct struct {
	StrVar   string
	intVar   int
	floatVar float64
	boolVar  bool
	child    ChildStruct
	childRef *ChildStruct
}

func PkgFuncVoid(strArg string) {
}

func pkgFuncStringReturn(strArg string) string {
	return "value"
}

func pkgFuncTupleReturn(strArg string) (string, error) {
	return "value", nil
}

func (p *ParentStruct) ParentPublicFunc(name string) (string, error) {
	return "value", nil
}

func (p *ParentStruct) parentPrivateFunc(name string) (string, error) {
	return "value", nil
}

func (p *ChildStruct) ChildPublicFunc(name string, num int) (string, error) {
	return "value", nil
}
