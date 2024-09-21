package model

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Convention represents a set of rules with a specific name.
type Convention struct {
	Name  string `json:"name"`
	Rules []Rule `json:"rules"`
}

// NewConvention creates a new Convention with the given name and rules.
func NewConvention(name string, rules []Rule) *Convention {
	return &Convention{Name: name, Rules: rules}
}

// Rule represents a single rule with a strictness level, type, and value.
type Rule struct {
	Level StrictnessLevel
	Type  RuleType
	Value string
}

// NewRule creates a new Rule with the given level, type, and value.
func NewRule(level StrictnessLevel, ruleType RuleType, value string) *Rule {
	return &Rule{Level: level, Type: ruleType, Value: value}
}

// MarshalJSON customizes the JSON representation of a Rule.
func (r *Rule) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"level": strictnessLevelNames[r.Level],
		"type":  ruleTypeNames[r.Type],
		"value": r.Value,
	})
}

// Print prints the convention's name and its rules.
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

// Evaluate performs an evaluation of the Convention.
func (c *Convention) Evaluate() error {
	return nil
}
