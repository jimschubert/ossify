package cmd

import (
	"github.com/spf13/cobra"
)

type CreateFlags struct {
	out string
}

func newCreateCmd() *cobra.Command {
	createFlags := &CreateFlags{}
	createCommand := &cobra.Command{
		Use:   "create",
		Short: "Create a new open-source project or other resource",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO
		},
	}

	licenseCmd := &cobra.Command{
		Use:   "license",
		Short: "Create or evaluate a new license for your repository.",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO
		},
	}

	projectCmd := &cobra.Command{
		Use:   "project",
		Short: "Create a new project.",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO
		},
	}

	repositoryCmd := &cobra.Command{
		Use:   "repository",
		Short: "Create a new remote repository.",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO
		},
	}

	createCommand.PersistentFlags().StringVarP(&createFlags.out, "out", "o", "",
		"Create a new open-source project or other resource in `dir`.")

	// new license
	createCommand.AddCommand(licenseCmd)

	// new project
	createCommand.AddCommand(projectCmd)

	// new repository
	createCommand.AddCommand(repositoryCmd)

	return createCommand
}

func init() {
	rootCmd.AddCommand(newCreateCmd())
}
