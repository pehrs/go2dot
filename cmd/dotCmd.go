package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"pehrs.com/go2dot/pkg/graphviz"
)

var dotCmd = &cobra.Command{
	Use:   "dot",
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Short: "Generate dot graph of structs and functions",
	Long:  `Generate a Graphviz dot file of your code`,
	RunE: func(cmd *cobra.Command, args []string) error {

		pkgPath := args[0]

		graphviz.ShowPrivate(showPrivate)

		fp, err := filepath.Abs(pkgPath)
		if err != nil {
			return err
		}
		info, err := os.Stat(fp)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return fmt.Errorf("Entry point should be a path to a directory of a go package.")
		}

		dot, err := graphviz.ToDot(pkgPath)
		if err != nil {
			return err
		}

		fmt.Printf("%s\n", dot)

		return nil
	},
}

// func toDot(pkgDir string) (string, error) {
// 	fset := token.NewFileSet()
// 	parsed, err := parser.ParseDir(fset, pkgDir, func(os.FileInfo) bool { return true }, parser.AllErrors)
// 	if err != nil {
// 		return "", err
// 	}
// 	dot := `
// 	digraph "golang" {
// 			fontname="Jetbrains Mono Regular,Ubuntu Mono,Helvetica"
// 			rankdir = TB;
// 			labelloc="t"
// 			graph [];
// 			node [
// 				fontname="Jetbrains Mono Regular,Ubuntu Mono,Helvetica"
// 				shape=record
// 				labelloc="t"
// 			];
// 	`
// 	for _, pkg := range parsed {

// 		var clusterCount = 0
// 		var fcount = 1
// 		for filename := range pkg.Files {
// 			structs, err := golang.ExtractStructs(filename)
// 			if err != nil {
// 				return "", err
// 			}
// 			var structMap = map[string](*golang.StructDecl){}
// 			var structNameMap = map[(*golang.StructDecl)]string{}
// 			for _, theStruct := range structs {
// 				if showPrivate || theStruct.IsPublic() {
// 					dot = dot + fmt.Sprintf("%s\n", theStruct.Dot(showPrivate))
// 					// typeId := theStruct.TypeId()
// 					structMap[theStruct.Name] = theStruct
// 					structNameMap[theStruct] = theStruct.Name
// 				}
// 			}
// 			for _, theStruct := range structs {
// 				if showPrivate || theStruct.IsPublic() {
// 					dot = dot + fmt.Sprintf("%s\n", theStruct.DotDeps(structNameMap, structMap))
// 				}
// 			}
// 		}

// 		dot = dot + fmt.Sprintf("subgraph cluster_%d { rank = same; label = \"«%s functions»\";\n",
// 			clusterCount, pkg.Name,
// 		)
// 		dot = dot + fmt.Sprintf("%s_Functions[label = <{",
// 			pkg.Name,
// 		)
// 		for filename := range pkg.Files {

// 			funcs, err := golang.ExtractFunctions(filename)
// 			if err != nil {
// 				return "", err
// 			}
// 			funcDecls := ""
// 			for _, theFunc := range funcs {
// 				if showPrivate || theFunc.IsPublic() {
// 					funcDecls = fmt.Sprintf("%s%s", funcDecls, theFunc.DotLabel())
// 				}
// 			}
// 			// fmt.Printf("%s_Functions_%d[label = <{<b>«func» (%s)</b><br align=\"left\"/>%s}>, shape=record];\n",
// 			// 	pkg.Name,
// 			// 	fcount,
// 			// 	filename,
// 			// 	funcDecls,
// 			// )
// 			dot = dot + fmt.Sprintf("<b>%s</b><br align=\"left\"/>%s", filename, funcDecls)
// 			fcount++
// 		}
// 		dot = dot + fmt.Sprintf("}>, color=white, shape=record];\n")
// 		dot = dot + fmt.Sprintf("}\n")
// 		clusterCount++

// 		dot = dot + fmt.Sprintf("}\n")
// 	}

// 	return dot, nil
// }

func init() {
	rootCmd.AddCommand(dotCmd)
	addConfigFlags(dotCmd)
}
