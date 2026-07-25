// Command athenaeum generates a podcast feed site from a library of .m4b
// audiobooks. The generated tree is static: serve it with any web server.
package main

import (
	"os"

	"github.com/spf13/cobra"
)

func main() {
	if err := NewRootCommand().Execute(); err != nil {
		os.Exit(1)
	}
}

// NewRootCommand builds the cobra command tree.
func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:          "athenaeum",
		Short:        "generates a podcast feed site from a library of audiobooks",
		SilenceUsage: true,
		Version:      Version,
	}

	rootCmd.AddCommand(NewVersionCommand())
	rootCmd.AddCommand(NewBuildCommand())
	return rootCmd
}
