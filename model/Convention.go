package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jimschubert/ossify/util"
	"strings"
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

//noinspection GoUnusedExportedFunction
func NewConvention(name string, rules []Rule) *Convention {
	return &Convention{Name: name, Rules: rules}
}

type Rule struct {
	Level StrictnessLevel
	Type  RuleType
	Value string
}

//noinspection GoUnusedExportedFunction
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
		return errors.New(fmt.Sprintf("type %s is not valid", other.Type))
	}

	level := util.StringSearch(strictnessLevelNames, other.Level)
	if level == -1 {
		return errors.New(fmt.Sprintf("level %s is not valid", other.Level))
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
	_, err := fmt.Printf(str.String())
	return err
}

func (c *Convention) Evaluate() error {
	return nil
}
