package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jimschubert/ossify/config"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:    "version",
	Short:  "Print the version number of ossify",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		// Version string can be tested with:
		// goreleaser release --skip-publish --snapshot --rm-dist
		var str strings.Builder
		str.WriteString(fmt.Sprintf("ossify %s (%s)\n", config.Version, config.Commit))
		str.WriteString(fmt.Sprintf("Built: %s", config.Date))
		fmt.Println(str.String())
	},
}
