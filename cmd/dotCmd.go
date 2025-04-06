package cmd

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"pehrs.com/go2dot/pkg/golang"
)

var showPrivate = false

var dotCmd = &cobra.Command{
	Use:   "dot",
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Short: "Generate dot graph of structs and functions",
	Long:  `Generate a Graphviz dot file of your codes`,
	RunE: func(cmd *cobra.Command, args []string) error {

		entry := args[0]

		fp, err := filepath.Abs(entry)
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
		fset := token.NewFileSet()
		parsed, err := parser.ParseDir(fset, entry, func(os.FileInfo) bool { return true }, parser.AllErrors)
		if err != nil {
			return err
		}
		for _, pkg := range parsed {
			fmt.Printf(`
	digraph "" {
			fontname="Jetbrains Mono Regular,Ubuntu Mono,Helvetica"
			rankdir = TB;
			labelloc="t"
			graph [];
			node [
				fontname="Jetbrains Mono Regular,Ubuntu Mono,Helvetica"
				shape=record
				labelloc="t"
			];
	`)

			var clusterCount = 0
			var fcount = 1
			for filename := range pkg.Files {
				structs, err := golang.ExtractStructs(filename)
				if err != nil {
					return err
				}
				var structMap = map[string](*golang.StructDecl){}
				var structNameMap = map[(*golang.StructDecl)]string{}
				for _, theStruct := range structs {
					if showPrivate || theStruct.IsPublic() {
						fmt.Printf("%s\n", theStruct.Dot(showPrivate))
						// typeId := theStruct.TypeId()
						structMap[theStruct.Name] = theStruct
						structNameMap[theStruct] = theStruct.Name
					}
				}
				for _, theStruct := range structs {
					if showPrivate || theStruct.IsPublic() {
						fmt.Printf("%s\n", theStruct.DotDeps(structNameMap, structMap))
					}
				}
			}

			fmt.Printf("subgraph cluster_%d { rank = same; label = \"«%s functions»\";\n",
				clusterCount, pkg.Name,
			)
			fmt.Printf("%s_Functions[label = <{",
				pkg.Name,
			)
			for filename := range pkg.Files {

				funcs, err := golang.ExtractFunctions(filename)
				if err != nil {
					return err
				}
				funcDecls := ""
				for _, theFunc := range funcs {
					if showPrivate || theFunc.IsPublic() {
						funcDecls = fmt.Sprintf("%s%s", funcDecls, theFunc.DotLabel())
					}
				}
				// fmt.Printf("%s_Functions_%d[label = <{<b>«func» (%s)</b><br align=\"left\"/>%s}>, shape=record];\n",
				// 	pkg.Name,
				// 	fcount,
				// 	filename,
				// 	funcDecls,
				// )
				fmt.Printf("<b>%s</b><br align=\"left\"/>%s", filename, funcDecls)
				fcount++
			}
			fmt.Printf("}>, color=white, shape=record];\n")
			fmt.Printf("}\n")
			clusterCount++

			fmt.Printf("}\n")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(dotCmd)

	dotCmd.PersistentFlags().BoolVarP(&showPrivate, "private", "p", showPrivate, "Render private structs and functions.")

}
