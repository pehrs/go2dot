package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"pehrs.com/go2dot/pkg/graphviz"
)

var additionalOptions = ""
var outputFormat = "png"

var graphCmd = &cobra.Command{
	Use:   "graph (pkg-path) (output-file)",
	Args:  cobra.MatchAll(cobra.ExactArgs(2), cobra.OnlyValidArgs),
	Short: "Generate dot graph of structs and functions",
	Long:  `Generate a image of your code`,
	RunE: func(cmd *cobra.Command, args []string) error {

		pkgPath := args[0]
		outputPath := args[1]

		graphviz.SetDotExec(dotExec)
		graphviz.SetOptions(additionalOptions)
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
		graphviz.RunDot(dot, outputFormat, outputPath)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(graphCmd)
	addConfigFlags(graphCmd)
	graphCmd.PersistentFlags().StringVarP(&additionalOptions, "options", "x", additionalOptions, "Additional options for dot command.")
	graphCmd.PersistentFlags().StringVarP(&outputFormat, "format", "T", outputFormat, "Output format.")

}
