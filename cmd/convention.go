package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/jimschubert/ossify/internal/config"
	"github.com/jimschubert/ossify/internal/config/conventions"
	"github.com/jimschubert/ossify/internal/model"
	"github.com/spf13/cobra"
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
	Use:   "add [file]",
	Args:  cobra.MaximumNArgs(1),
	Short: "Adds a new custom convention (local-only) to the list of known conventions",
	Long: `Adds a new custom convention to your local configuration.

The convention can be provided as a JSON file path, or piped via stdin.
If --id is provided, it will be used as the filename; otherwise the convention's name is used.

Example JSON format:
{
  "name": "My Convention",
  "rules": [
    { "level": "required", "type": "directory", "value": "src" },
    { "level": "optional", "type": "file", "value": "CONTRIBUTING.md" }
  ]
}

Valid levels: prohibited, optional, preferred, required
Valid types: directory, file, pattern`,
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.ConfigManager.Load()
		failOnError(err)

		conventionsPath := conf.ConventionPath
		if conventionsPath == "" {
			fmt.Println("invalid conventions path: please update your configuration and try again")
			os.Exit(1)
		}

		err = os.MkdirAll(conventionsPath, 0755)
		failOnError(err)

		var data []byte

		if len(args) == 1 {
			// Read from file
			filePath := args[0]
			data, err = os.ReadFile(filePath)
			if err != nil {
				fmt.Printf("failed to read file %s: %v\n", filePath, err)
				os.Exit(1)
			}
		} else {
			// Read from stdin
			stat, _ := os.Stdin.Stat()
			if (stat.Mode() & os.ModeCharDevice) != 0 {
				fmt.Println("no input provided: specify a file path or pipe JSON via stdin")
				_ = cmd.Help()
				os.Exit(1)
			}
			data, err = io.ReadAll(os.Stdin)
			failOnError(err)
		}

		// Parse and validate the convention
		var convention model.Convention
		err = json.Unmarshal(data, &convention)
		if err != nil {
			fmt.Printf("invalid convention JSON: %v\n", err)
			os.Exit(1)
		}

		if convention.Name == "" {
			fmt.Println("convention must have a name")
			os.Exit(1)
		}

		if len(convention.Rules) == 0 {
			fmt.Println("convention must have at least one rule")
			os.Exit(1)
		}

		// Determine the filename
		id := conventionFlags.id
		if id == "" {
			id = convention.Name
		}

		// Sanitize the id for use as a filename
		id = strings.ReplaceAll(id, " ", "-")
		id = strings.ReplaceAll(id, "/", "-")
		id = strings.ToLower(id)

		filename := filepath.Join(conventionsPath, id+".json")

		// Check if file already exists
		if _, err := os.Stat(filename); err == nil {
			fmt.Printf("convention '%s' already exists at %s\n", id, filename)
			fmt.Println("use a different --id or remove the existing file")
			os.Exit(1)
		}

		// Marshal with indentation for readability
		output, err := json.MarshalIndent(convention, "", "  ")
		failOnError(err)

		err = os.WriteFile(filename, output, 0644)
		failOnError(err)

		fmt.Printf("convention '%s' saved to %s\n", convention.Name, filename)
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
