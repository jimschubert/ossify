package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jimschubert/ossify/config/conventions"
)

var conventionFlags *ConventionFlags

type ConventionFlags struct {
	id string
}

func init() {
	conventionFlags = &ConventionFlags{}
	rootCmd.AddCommand(conventionCmd)

	// convention
	conventionCmd.AddCommand(addConventionCmd)
	conventionCmd.AddCommand(listConventionCmd)

	// convention add
	addConventionCmd.Flags().StringVarP(&conventionFlags.id, "id", "i", "",
		"The identifier to be associated with your customized convention. This will take precedence over a built-in convention with the same id.")
}

var conventionCmd = &cobra.Command{
	Use:   "convention",
	Short: "Manage file and structure conventions",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO
	},
}

var addConventionCmd = &cobra.Command{
	Use:   "add [options]",
	Args:  cobra.MinimumNArgs(1),
	Short: "Adds a new custom convention (local-only) to the list of known conventions",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO
		// config, err := config.LoadConfig()
		// failOnError(err)
		//
		// conventionsPath := config.ConventionsPath
		// if conventionsPath == "" {
		//	fmt.Println("invalid conventions path: please update your configuration and try again")
		//  os.Exit(1)
		// }
		//
		// err = os.MkdirAll(conventionsPath, 0700)
		// failOnError(err)
	},
}

var listConventionCmd = &cobra.Command{
	Use:   "list",
	Short: "Presents a list of known conventions.",
	Run: func(cmd *cobra.Command, args []string) {
		conventions, err := conventions.Load()
		failOnError(err)

		for _, c := range *conventions {
			err = c.Print()
			if err != nil {
				break
			}
		}
	},
}
