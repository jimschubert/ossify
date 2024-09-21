package model

import (
	"encoding/json"
	"fmt"
	"slices"
)

// RuleType represents the type of a rule.
type RuleType int

// String returns the string representation of a RuleType.
func (r RuleType) String() string {
	switch r {
	case Directory:
		return "directory"
	case File:
		return "file"
	case Pattern:
		return "pattern"
	case FilesystemPolicy:
		return "filesystem policy"
	default:
		return "unspecified"
	}
}

// UnmarshalJSON customizes the JSON deserialization of a RuleType.
func (r *RuleType) UnmarshalJSON(data []byte) error {
	var ruleType int
	var ruleTypeName string
	if err := json.Unmarshal(data, &ruleTypeName); err == nil {
		ruleType = slices.Index(ruleTypeNames, ruleTypeName)
	} else if err := json.Unmarshal(data, &ruleType); err != nil {
		return err
	}

	newRule := RuleType(ruleType)

	if newRule < Unspecified || newRule >= FilesystemPolicy {
		return fmt.Errorf("rule type %d is not valid", ruleType)
	}

	*r = newRule
	return nil
}

const (
	Unspecified RuleType = iota
	Directory
	File
	Pattern
	FilesystemPolicy
)

var ruleTypeNames = []string{
	Unspecified.String(),
	Directory.String(),
	File.String(),
	Pattern.String(),
	FilesystemPolicy.String(),
}
