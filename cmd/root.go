package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var version bool
var rootCmd = &cobra.Command{
	Use:   "ossify",
	Short: "Give some structure to your open-source software projects.",
	Long: "Template, evaluate, and bootstrap files and directory for your open source projects.\n" +
		"Complete documentation is available at https://github.com/jimschubert/ossify.\n\n" +
		"\"Give it some bones!\"",
	Run: func(cmd *cobra.Command, args []string) {
		if version {
			versionCmd.Run(cmd, args)
			os.Exit(0)
		} else {
			_ = cmd.Help()
			os.Exit(1)
		}
	},
}

func failOnError(err error){
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

func init(){
	rootCmd.Flags().BoolVarP(&version, "version", "v", false, "version for ossify")
}