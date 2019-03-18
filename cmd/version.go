package cmd

import (
	"fmt"
	"github.com/jimschubert/ossify/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of ossify",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: consider using goxc to set version on build.
		// see https://sbstjn.com/create-golang-cli-application-with-cobra-and-goxc.html
		fmt.Printf("ossify %s\n", config.Version)
	},
}