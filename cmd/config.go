package cmd

import (
	"github.com/spf13/cobra"
)

var showPrivate = false
var dotExec = "dot"
var verbose = false

func addConfigFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&showPrivate, "private", "p", showPrivate, "Render private structs and functions.")
	cmd.PersistentFlags().StringVarP(&dotExec, "dot", "d", dotExec, "Graphviz DOT executable to use for rendering images.")
	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", verbose, "Be verbose")
}
