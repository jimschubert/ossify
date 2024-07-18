package model

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

type License struct {
	SupersededBy *string      `json:"superseded_by"`
	Identifiers  []Identifier `json:"identifiers"`
	Text         []Text       `json:"text"`
	OtherNames   *[]OtherName `json:"other_names"`
	Id           string       `json:"id"`
	Links        []Link       `json:"links"`
	Name         string       `json:"name"`
	Keywords     []string     `json:"keywords"`
}

// This type alias allows us to create a sort of "dao" atop the set of data.
// see Go in Action Chapter 2 for somewhat similar approach.
type Licenses []License

func (license License) Print() error {
	d := color.New(color.FgWhite, color.Bold)
	_, err := d.Printf("%-20s(%s)\n", license.Id, license.Name)
	return err
}

func (license License) PrintDetails() error {
	builder := strings.Builder{}

	bold := color.New(color.FgWhite, color.Bold)
	italic := color.New(color.FgWhite, color.Italic)
	warn := color.New(color.FgRed, color.Italic)
	builder.WriteString(bold.Sprintf("%-20s(%s)\n", license.Id, license.Name))
	if len(license.Keywords) > 0 {
		builder.WriteString(italic.Sprint(strings.Join(license.Keywords, ", ")))
		builder.WriteString("\n")
	}
	if license.SupersededBy != nil {
		builder.WriteString(warn.Sprintf("This license is superseded by %s\n", *license.SupersededBy))
	}

	if license.OtherNames != nil && len(*license.OtherNames) > 0 {
		builder.WriteString(bold.Sprintln("\nCommon names"))
		for _, other := range *license.OtherNames {
			builder.WriteString(fmt.Sprintf("  * %s\n", other.Name))
		}
	}

	if len(license.Identifiers) > 0 {
		builder.WriteString(bold.Sprintln("\nLicense Standards"))
		for _, identifier := range license.Identifiers {
			builder.WriteString(fmt.Sprintf("  * %-10s %s\t\n", identifier.Scheme, identifier.Identifier))
		}
	}

	if len(license.Links) > 0 {
		builder.WriteString(bold.Sprintln("\nLinks"))
		for _, link := range license.Links {
			builder.WriteString("  * ")
			builder.WriteString(link.Url)
			if link.Note != nil {
				builder.WriteString(italic.Sprintf(" (%s)", *link.Note))
			}
			builder.WriteString("\n")
		}
	}

	_, err := fmt.Print(builder.String())
	return err
}

// enable searching for a specific license
func (l Licenses) FindById(id string) *License {
	for _, license := range l {
		if strings.Contains(strings.ToLower(license.Id), strings.ToLower(id)) {
			return &license
		}
	}
	return nil
}

// enable querying by keyword (e.g. "popular")
func (l Licenses) FindByKeyword(keyword string) *Licenses {
	var licenses Licenses
	for _, license := range l {
		for _, kw := range license.Keywords {
			if kw == keyword {
				licenses = append(licenses, license)
			}
		}
	}
	return &licenses
}

// Search enables a loose free-form textual search
func (l Licenses) Search(term string) *Licenses {
	var licenses Licenses
SearchLoop:
	for _, license := range l {
		// lowercase compare
		if strings.Contains(strings.ToLower(license.Id), strings.ToLower(term)) {
			licenses = append(licenses, license)
			continue SearchLoop
		}

		if license.Identifiers != nil {
			for _, identifier := range license.Identifiers {
				if "SPDX" == identifier.Scheme {
					// use exact case-sensitive match
					if term == identifier.Identifier {
						licenses = append(licenses, license)
						continue SearchLoop
					}
				} else if fuzzy.MatchNormalizedFold(term, identifier.Identifier) {
					licenses = append(licenses, license)
					continue SearchLoop
				}
			}
		}

		if fuzzy.MatchNormalizedFold(term, license.Name) {
			licenses = append(licenses, license)
			continue SearchLoop
		}

		if license.OtherNames != nil {
			for _, otherNames := range *license.OtherNames {
				// We don't really care about the Note field here because it's just supplemental to the name.
				if fuzzy.MatchNormalizedFold(term, otherNames.Name) {
					licenses = append(licenses, license)
					continue SearchLoop
				}
			}
		}

		// TODO: search license contents
	}
	return &licenses
}
