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

		graphviz.Verbose(verbose)
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

func init() {
	rootCmd.AddCommand(dotCmd)
	addConfigFlags(dotCmd)
}
