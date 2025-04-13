package cmd

import (
	"github.com/spf13/cobra"
)

var showPrivate = false

func addConfigFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&showPrivate, "private", "p", showPrivate, "Render private structs and functions.")
}
