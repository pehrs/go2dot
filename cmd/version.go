package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "n/a"
	date    = "n/a"
	verbose = false
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show Version",
	RunE: func(cmd *cobra.Command, args []string) error {
		if verbose {
			fmt.Printf("%s, %s, %s\n", version, commit, date)
		} else {
			fmt.Printf("%s\n", version)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", verbose, "Be verbose")

}
