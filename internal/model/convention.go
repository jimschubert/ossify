package model

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jimschubert/ossify/internal/util"
)

type StrictnessLevel int
type RuleType int

const (
	Prohibited StrictnessLevel = iota
	Optional
	Preferred
	Required
)

const (
	Unspecified RuleType = iota
	Directory
	File
	Pattern
)

var strictnessLevelNames = []string{
	Prohibited: "prohibited",
	Optional:   "optional",
	Preferred:  "preferred",
	Required:   "required",
}

var ruleTypeNames = []string{
	Unspecified: "unspecified",
	Directory:   "directory",
	File:        "file",
	Pattern:     "pattern",
}

type Convention struct {
	Name  string `json:"name"`
	Rules []Rule `json:"rules"`
}

// noinspection GoUnusedExportedFunction
func NewConvention(name string, rules []Rule) *Convention {
	return &Convention{Name: name, Rules: rules}
}

type Rule struct {
	Level StrictnessLevel
	Type  RuleType
	Value string
}

// noinspection GoUnusedExportedFunction
func NewRule(level StrictnessLevel, ruleType RuleType, value string) *Rule {
	return &Rule{Level: level, Type: ruleType, Value: value}
}

func (r *Rule) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"level": strictnessLevelNames[r.Level],
		"type":  ruleTypeNames[r.Type],
		"value": r.Value,
	})
}

func (r *Rule) UnmarshalJSON(data []byte) error {
	other := &struct {
		Level string `json:"level"`
		Type  string `json:"type"`
		Value string `json:"value"`
	}{}

	if err := json.Unmarshal(data, &other); err != nil {
		return err
	}

	ruleType := util.StringSearch(ruleTypeNames, other.Type)
	if ruleType == -1 {
		return fmt.Errorf("type %s is not valid", other.Type)
	}

	level := util.StringSearch(strictnessLevelNames, other.Level)
	if level == -1 {
		return fmt.Errorf("level %s is not valid", other.Level)
	}

	r.Value = other.Value
	r.Type = RuleType(ruleType)
	r.Level = StrictnessLevel(level)

	return nil
}

func (c *Convention) Print() error {
	var str strings.Builder
	str.WriteString(c.Name)
	if len(c.Rules) == 0 {
		str.WriteString(": No Rules Specified!\n")
	} else {
		str.WriteString("\n")
		for _, r := range c.Rules {
			str.WriteString(fmt.Sprintf("  - %-20s %-15s %-10s\n", r.Value, ruleTypeNames[r.Type], strictnessLevelNames[r.Level]))
		}
	}
	_, err := fmt.Print(str.String())
	return err
}

// RuleResult represents the result of evaluating a single rule
type RuleResult struct {
	Rule    Rule
	Passed  bool
	Message string
}

// CheckResult represents the overall result of checking a convention against a directory
type CheckResult struct {
	Convention  string
	Directory   string
	Results     []RuleResult
	PassCount   int
	FailCount   int
	WarnCount   int
	SkipCount   int
}

// HasFailures returns true if any required rules failed or prohibited items exist
func (cr *CheckResult) HasFailures() bool {
	return cr.FailCount > 0
}

// Print outputs the check results to stdout
func (cr *CheckResult) Print() {
	fmt.Printf("Checking convention '%s' against directory: %s\n\n", cr.Convention, cr.Directory)

	for _, r := range cr.Results {
		status := "✓"
		if !r.Passed {
			switch r.Rule.Level {
			case Required, Prohibited:
				status = "✗"
			case Preferred:
				status = "⚠"
			default:
				status = "○"
			}
		}
		fmt.Printf("  %s %-20s %-12s %-10s %s\n",
			status,
			r.Rule.Value,
			ruleTypeNames[r.Rule.Type],
			strictnessLevelNames[r.Rule.Level],
			r.Message)
	}

	fmt.Printf("\nSummary: %d passed, %d failed, %d warnings, %d skipped\n",
		cr.PassCount, cr.FailCount, cr.WarnCount, cr.SkipCount)
}

// Evaluate checks all rules in the convention against the specified directory
func (c *Convention) Evaluate(targetDir string) (*CheckResult, error) {
	result := &CheckResult{
		Convention: c.Name,
		Directory:  targetDir,
		Results:    make([]RuleResult, 0, len(c.Rules)),
	}

	for _, rule := range c.Rules {
		ruleResult := evaluateRule(rule, targetDir)
		result.Results = append(result.Results, ruleResult)

		if ruleResult.Passed {
			result.PassCount++
		} else {
			switch rule.Level {
			case Required, Prohibited:
				result.FailCount++
			case Preferred:
				result.WarnCount++
			case Optional:
				result.SkipCount++
			}
		}
	}

	return result, nil
}

func evaluateRule(rule Rule, targetDir string) RuleResult {
	result := RuleResult{Rule: rule}

	targetPath := filepath.Join(targetDir, rule.Value)
	info, err := os.Stat(targetPath)
	exists := err == nil

	switch rule.Type {
	case Directory:
		if exists && info.IsDir() {
			// Directory exists
			if rule.Level == Prohibited {
				result.Passed = false
				result.Message = "prohibited directory exists"
			} else {
				result.Passed = true
				result.Message = "found"
			}
		} else {
			// Directory does not exist
			if rule.Level == Prohibited {
				result.Passed = true
				result.Message = "not present (good)"
			} else if rule.Level == Required {
				result.Passed = false
				result.Message = "missing"
			} else if rule.Level == Preferred {
				result.Passed = false
				result.Message = "recommended but missing"
			} else {
				result.Passed = true
				result.Message = "not present (optional)"
			}
		}

	case File:
		if exists && !info.IsDir() {
			// File exists
			if rule.Level == Prohibited {
				result.Passed = false
				result.Message = "prohibited file exists"
			} else {
				result.Passed = true
				result.Message = "found"
			}
		} else {
			// File does not exist
			if rule.Level == Prohibited {
				result.Passed = true
				result.Message = "not present (good)"
			} else if rule.Level == Required {
				result.Passed = false
				result.Message = "missing"
			} else if rule.Level == Preferred {
				result.Passed = false
				result.Message = "recommended but missing"
			} else {
				result.Passed = true
				result.Message = "not present (optional)"
			}
		}

	case Pattern:
		// Pattern matching using glob
		matches, globErr := filepath.Glob(filepath.Join(targetDir, rule.Value))
		if globErr != nil {
			result.Passed = false
			result.Message = fmt.Sprintf("invalid pattern: %v", globErr)
		} else if len(matches) > 0 {
			if rule.Level == Prohibited {
				result.Passed = false
				result.Message = fmt.Sprintf("prohibited pattern matched %d item(s)", len(matches))
			} else {
				result.Passed = true
				result.Message = fmt.Sprintf("matched %d item(s)", len(matches))
			}
		} else {
			if rule.Level == Prohibited {
				result.Passed = true
				result.Message = "no matches (good)"
			} else if rule.Level == Required {
				result.Passed = false
				result.Message = "no matches"
			} else if rule.Level == Preferred {
				result.Passed = false
				result.Message = "recommended but no matches"
			} else {
				result.Passed = true
				result.Message = "no matches (optional)"
			}
		}

	default:
		result.Passed = false
		result.Message = "unknown rule type"
	}

	return result
}
