package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jimschubert/ossify/internal/config/conventions"
	"github.com/jimschubert/ossify/internal/model"
	"github.com/spf13/cobra"
)

var checkFlags *CheckFlags

// CheckFlags holds the flag values for the check command
type CheckFlags struct {
	conventionID   string
	conventionFile string
	directory      string
	all            bool
}

func init() {
	checkFlags = &CheckFlags{}
	rootCmd.AddCommand(checkCmd)

	checkCmd.Flags().StringVarP(&checkFlags.conventionID, "convention", "c", "",
		"The ID/name of a convention to check against (e.g., 'Go', 'Standard Distribution')")
	checkCmd.Flags().StringVarP(&checkFlags.conventionFile, "file", "f", "",
		"Path to a JSON file describing the convention rules to check")
	checkCmd.Flags().StringVarP(&checkFlags.directory, "directory", "d", ".",
		"The directory to check (defaults to current directory)")
	checkCmd.Flags().BoolVarP(&checkFlags.all, "all", "a", false,
		"Check against all known conventions")
}

var checkCmd = &cobra.Command{
	Use:   "check [convention-name]",
	Short: "Check a directory against convention rules",
	Long: `Check a directory against one or more convention rules.

You can specify the convention to check against in several ways:
  1. By name as an argument: ossify check Go
  2. By name with --convention flag: ossify check -c "Standard Distribution"
  3. By JSON file: ossify check -f my-convention.json
  4. Check all conventions: ossify check --all

The directory to check defaults to the current directory, but can be
specified with the --directory flag.

Exit codes:
  0 - All required rules pass
  1 - One or more required rules failed`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check for mutually exclusive options
		conventionID := checkFlags.conventionID
		if conventionID == "" && len(args) > 0 {
			conventionID = args[0]
		}

		optionsCount := 0
		if checkFlags.conventionFile != "" {
			optionsCount++
		}
		if conventionID != "" {
			optionsCount++
		}
		if checkFlags.all {
			optionsCount++
		}

		if optionsCount > 1 {
			cobra.CheckErr(fmt.Errorf("--file, --convention (or convention name argument), and --all are mutually exclusive"))
		}

		// Determine the target directory
		targetDir := checkFlags.directory
		if targetDir == "" {
			targetDir = "."
		}

		// Convert to absolute path
		absDir, err := filepath.Abs(targetDir)
		if err != nil {
			cobra.CheckErr(fmt.Errorf("resolving directory path: %w", err))
		}

		// Verify directory exists
		info, err := os.Stat(absDir)
		if err != nil {
			cobra.CheckErr(fmt.Errorf("accessing directory: %w", err))
		}
		if !info.IsDir() {
			cobra.CheckErr(fmt.Errorf("'%s' is not a directory", absDir))
		}

		// Collect conventions to check
		var conventionsToCheck []model.Convention

		// Option 1: Load from JSON file
		if checkFlags.conventionFile != "" {
			convention, err := loadConventionFromFile(checkFlags.conventionFile)
			if err != nil {
				cobra.CheckErr(fmt.Errorf("loading convention file: %w", err))
			}
			conventionsToCheck = append(conventionsToCheck, *convention)
		}

		// Option 2: Find by ID/name (from flag or argument)
		if conventionID != "" {
			allConventions, err := conventions.Load()
			if err != nil {
				cobra.CheckErr(fmt.Errorf("loading conventions: %w", err))
			}

			convention := findConventionByName(*allConventions, conventionID)
			if convention == nil {
				fmt.Printf("convention '%s' not found\n", conventionID)
				fmt.Println("\nAvailable conventions:")
				for _, c := range *allConventions {
					fmt.Printf("  - %s\n", c.Name)
				}
				os.Exit(1)
			}
			conventionsToCheck = append(conventionsToCheck, *convention)
		}

		// Option 3: Check all conventions
		if checkFlags.all {
			allConventions, err := conventions.Load()
			if err != nil {
				cobra.CheckErr(fmt.Errorf("loading conventions: %w", err))
			}
			conventionsToCheck = *allConventions
		}

		// If no convention specified, show help
		if len(conventionsToCheck) == 0 {
			fmt.Println("no convention specified")
			fmt.Println()
			_ = cmd.Help()
			os.Exit(1)
		}

		// Run checks
		hasFailures := false
		for i, convention := range conventionsToCheck {
			if i > 0 {
				fmt.Println("\n" + strings.Repeat("-", separatorWidth) + "\n")
			}

			result, err := convention.Evaluate(absDir)
			if err != nil {
				cobra.CheckErr(fmt.Errorf("evaluating convention '%s': %w", convention.Name, err))
			}

			result.Print()

			if result.HasFailures() {
				hasFailures = true
			}
		}

		if hasFailures {
			os.Exit(1)
		}
	},
}

const separatorWidth = 60

// loadConventionFromFile loads and validates a convention from a JSON file.
// If the convention has no name, the filename is used as the name.
func loadConventionFromFile(filePath string) (*model.Convention, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var convention model.Convention
	if err := json.Unmarshal(data, &convention); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	if convention.Name == "" {
		convention.Name = filepath.Base(filePath)
	}

	if len(convention.Rules) == 0 {
		return nil, fmt.Errorf("convention must have at least one rule")
	}

	return &convention, nil
}

// findConventionByName searches for a convention by name (case-insensitive).
// Returns nil if no matching convention is found.
func findConventionByName(conventions []model.Convention, name string) *model.Convention {
	nameLower := strings.ToLower(name)
	for _, c := range conventions {
		if strings.ToLower(c.Name) == nameLower {
			return &c
		}
	}
	return nil
}
