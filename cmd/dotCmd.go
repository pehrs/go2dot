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
			title := pkg.Name
			fmt.Printf(`
	digraph "" {
			labelloc="t"
			graph [label = "%s"];
			node [
					shape=record
					labelloc="t"
			];
	`, title)

			for filename := range pkg.Files {
				structs, err := golang.ExtractStructs(filename)
				if err != nil {
					return err
				}
				var structMap = map[string](*golang.GolangStruct){}
				var structNameMap = map[(*golang.GolangStruct)]string{}
				for _, theStruct := range structs {
					fmt.Printf("%s\n", theStruct.Dot())
					// typeId := theStruct.TypeId()
					structMap[theStruct.Name] = theStruct
					structNameMap[theStruct] = theStruct.Name
				}
				for _, theStruct := range structs {
					fmt.Printf("%s\n", theStruct.DotDeps(structNameMap, structMap))
				}
			}
			fmt.Printf("}\n")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(dotCmd)

}
