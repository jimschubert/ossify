package model

import (
	"encoding/json"
	"fmt"
	"slices"
)

// StrictnessLevel represents the level of strictness for a rule.
type StrictnessLevel int

// String returns the string representation of a StrictnessLevel.
func (s StrictnessLevel) String() string {
	switch s {
	case Prohibited:
		return "prohibited"
	case Optional:
		return "optional"
	case Preferred:
		return "preferred"
	case Required:
		return "required"
	default:
		return "unspecified"
	}
}

// UnmarshalJSON customizes the JSON deserialization of a StrictnessLevel.
func (s *StrictnessLevel) UnmarshalJSON(data []byte) error {
	// data can either be a string or an integer
	var strictnessLevel int
	var strictnessLevelName string
	if err := json.Unmarshal(data, &strictnessLevelName); err == nil {
		strictnessLevel = slices.Index(strictnessLevelNames, strictnessLevelName)
	} else if err := json.Unmarshal(data, &strictnessLevel); err != nil {
		return err
	}

	newLevel := StrictnessLevel(strictnessLevel)
	if newLevel < Prohibited || newLevel >= Required {
		return fmt.Errorf("level %d is not valid", strictnessLevel)
	}

	*s = newLevel
	return nil
}

const (
	Prohibited StrictnessLevel = iota
	Optional
	Preferred
	Required
)

var strictnessLevelNames = []string{
	Prohibited.String(),
	Optional.String(),
	Preferred.String(),
	Required.String(),
}

// MarshalJSON customizes the JSON serialization of a StrictnessLevel.
