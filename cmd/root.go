package cmd

import (
	"fmt"
	"github.com/jimschubert/ossify/internal/config"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = newRootCmd()

func newRootCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "ossify",
		Short: "Give some structure to your open-source software projects.",
		Long: "Template, evaluate, and bootstrap files and directory for your open source projects.\n" +
			"Complete documentation is available at https://github.com/jimschubert/ossify.\n\n" +
			"\"Give it some bones!\"",
		Version: fmt.Sprintf("ossify %s (%s)\nBuilt: %s", config.Version, config.Commit, config.Date),
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
			os.Exit(1)
		},
	}

	return c
}

func failOnError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
