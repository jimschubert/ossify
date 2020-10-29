package cmd

import (
	"github.com/spf13/cobra"
)

var newFlags *NewFlags

type NewFlags struct {
	out string
}

func init() {
	newFlags = &NewFlags{}
	rootCmd.AddCommand(newCmd)

	newCmd.PersistentFlags().StringVarP(&newFlags.out, "out", "o", "",
		"Create a new open-source project or other resource in `dir`.")

	// new license
	newCmd.AddCommand(newLicenseCmd)

	// new project
	newCmd.AddCommand(newProjectCmd)

	// new repository
	newCmd.AddCommand(newRepositoryCmd)
}

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new open-source project or other resource",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO
	},
}

var newLicenseCmd = &cobra.Command{
	Use:   "license",
	Short: "Create or evaluate a new license for your repository.",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO
	},
}

var newProjectCmd = &cobra.Command{
	Use:   "project",
	Short: "Create a new project.",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO
	},
}

var newRepositoryCmd = &cobra.Command{
	Use:   "repository",
	Short: "Create a new remote repository.",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO
	},
}
