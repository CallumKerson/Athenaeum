package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version = "development"
	Commit  = "development"
	Date    = "development"
)

func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:               "version",
		Short:             "Print version and build information",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error { return nil },
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(cmd.OutOrStdout(), "version: ", Version)
			fmt.Fprintln(cmd.OutOrStdout(), "commit:  ", Commit)
			fmt.Fprintln(cmd.OutOrStdout(), "built at:", Date)
		},
	}
}
