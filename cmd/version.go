package cmd

import (
	"fmt"
	"github.com/jimschubert/ossify/config"
	"github.com/spf13/cobra"
	"strings"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of ossify",
	Run: func(cmd *cobra.Command, args []string) {
		// Version string can be tested with:
		// goreleaser release --skip-publish --snapshot --rm-dist
		var str strings.Builder
		str.WriteString(fmt.Sprintf("ossify %s (%s)\n", config.Version, config.Commit))
		str.WriteString(fmt.Sprintf("Built: %s", config.Date))
		fmt.Println(str.String())
	},
}