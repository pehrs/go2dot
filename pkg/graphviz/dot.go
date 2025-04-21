package graphviz

import (
	"bufio"
	"bytes"
	"fmt"
	"go/parser"
	"go/token"
	"log"
	"os"
	"os/exec"

	"pehrs.com/go2dot/pkg/golang"
)

var dotExec = "dot"
var showPrivate = false
var verbose = false

func Verbose(v bool) {
	verbose = v
}

// More Configs
// -Gfontname="Sans" -Nfontname="Sans" -Gsize=4,3 -Gdpi=1000

var options = ""

func ShowPrivate(value bool) {
	showPrivate = value
}

func SetDotExec(value string) {
	dotExec = value
}

func SetOptions(opts string) {
	options = opts
}

var clusterCount = 0

func ToDot(pkgDir string) (string, error) {
	fset := token.NewFileSet()
	parsed, err := parser.ParseDir(fset, pkgDir, func(os.FileInfo) bool { return true }, parser.AllErrors)
	if err != nil {
		return "", err
	}
	dot := `
	digraph "golang" {
			fontname="Jetbrains Mono Regular,Ubuntu Mono,Helvetica"
			rankdir = TB;
			labelloc="t"
			graph [];
			node [
				fontname="Jetbrains Mono Regular,Ubuntu Mono,Helvetica"
				shape=record
				labelloc="t"
			];
	`
	for _, pkg := range parsed {

		var fcount = 1
		for filename := range pkg.Files {
			structs, err := golang.ExtractStructs(filename)
			if err != nil {
				return "", err
			}
			var structMap = map[string](*golang.StructDecl){}
			var structNameMap = map[(*golang.StructDecl)]string{}
			for _, theStruct := range structs {
				if showPrivate || theStruct.IsPublic() {
					dot = dot + fmt.Sprintf("%s\n", theStruct.Dot(showPrivate))
					// typeId := theStruct.TypeId()
					structMap[theStruct.Name] = theStruct
					structNameMap[theStruct] = theStruct.Name
				}
			}
			for _, theStruct := range structs {
				if showPrivate || theStruct.IsPublic() {
					dot = dot + fmt.Sprintf("%s\n", theStruct.DotDeps(structNameMap, structMap))
				}
			}
		}

		dot = dot + fmt.Sprintf("subgraph cluster_%d { rank = same; label = \"«pkg:%s»\";\n",
			clusterCount, pkg.Name,
		)
		clusterCount++

		dot = dot + fmt.Sprintf("%s_Functions[label = <{",
			pkg.Name,
		)
		for filename := range pkg.Files {

			vars, err := golang.ExtractVars(filename)
			if err != nil {
				return "", err
			}
			varDecls := ""
			for _, theVar := range vars {
				if showPrivate || theVar.IsPublic() {
					varDecls = fmt.Sprintf("%s%s", varDecls, theVar.DotLabel())
				}
			}

			funcs, err := golang.ExtractFunctions(filename)
			if err != nil {
				return "", err
			}
			funcDecls := ""
			for _, theFunc := range funcs {
				if showPrivate || theFunc.IsPublic() {
					funcDecls = fmt.Sprintf("%s%s", funcDecls, theFunc.DotLabel())
				}
			}
			if fcount > 1 {
				dot = dot + "<br align=\"left\"/>"
			}
			dot = dot + fmt.Sprintf("<b>%s</b><br align=\"left\"/>%s<br align=\"left\"/>%s", filename, varDecls, funcDecls)
			fcount++
		}
		dot = dot + fmt.Sprintf("}>, color=white, shape=record];\n")

		dot = dot + fmt.Sprintf("}\n")
	}
	dot = dot + fmt.Sprintf("}\n")

	return dot, nil
}

func RunDot(dotGraph, format, outputFilename string) error {

	dotFile, err := os.CreateTemp("/tmp", "go2dot")
	if err != nil {
		log.Fatal(err)
	}

	w := bufio.NewWriter(dotFile)
	w.WriteString(dotGraph)
	w.Flush()
	dotFile.Close()

	err = runDotForFile(dotFile.Name(), format, outputFilename)
	defer os.Remove(dotFile.Name())
	if err != nil {
		return err
	}
	return nil
}

func runDotForFile(dotFilename, format, outputFilename string) error {

	bashCmdLine := fmt.Sprintf("%s %s -T%s -o %s", dotExec, dotFilename, format, outputFilename)
	if verbose {
		fmt.Printf("comand line:\n%v\n", bashCmdLine)
	}
	if len(options) > 0 {
		// optParts := strings.Split(options, " ")
		// optProcessedParts := []string{}
		// for _, part := range optParts {
		// 	optProcessedParts = append(optProcessedParts, strings.Replace(part, "\"", "", -1))
		// }
		bashCmdLine = fmt.Sprintf("%s %s", bashCmdLine, options)
	}

	// cmd := exec.Command(dotExec, cmdLine...)
	cmd := exec.Command("bash", "-c", bashCmdLine)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	fmt.Println(outb.String(), errb.String())
	if err != nil {
		fmt.Printf("ERROR WHEN CALLING DOT! %v\n\n", err)
		return err
	}
	return nil
}
