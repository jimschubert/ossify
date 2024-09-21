package model

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

// License represents a software license with various attributes.
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

// The Licenses type alias allows us to create a sort of "dao" atop the set of data.
// see Go in Action Chapter 2 for somewhat similar approach.
type Licenses []License

func (license License) Print() error {
	d := color.New(color.FgWhite, color.Bold)
	_, err := d.Printf("%-20s(%s)\n", license.Id, license.Name)
	return err
}

// PrintDetails prints the details of a license to the console.
func (license License) PrintDetails() error {
	builder := strings.Builder{}

	bold := color.New(color.FgWhite, color.Bold)
	italic := color.New(color.FgWhite, color.Italic)
	warn := color.New(color.FgRed, color.Italic)

	builder.WriteString(bold.Sprintf("%-20s(%s)\n", license.Id, license.Name))
	license.appendKeywords(&builder, italic)
	license.appendSupersededBy(&builder, warn)
	license.appendOtherNames(&builder, bold)
	license.appendIdentifiers(&builder, bold)
	license.appendLinks(&builder, bold, italic)

	_, err := fmt.Print(builder.String())
	return err
}

// appendKeywords appends the keywords to the license details.
func (license License) appendKeywords(builder *strings.Builder, italic *color.Color) {
	if len(license.Keywords) > 0 {
		builder.WriteString(italic.Sprint(strings.Join(license.Keywords, ", ")))
		builder.WriteString("\n")
	}
}

// appendSupersededBy appends the superseded by information to the license details.
func (license License) appendSupersededBy(builder *strings.Builder, warn *color.Color) {
	if license.SupersededBy != nil {
		builder.WriteString(warn.Sprintf("This license is superseded by %s\n", *license.SupersededBy))
	}
}

// appendOtherNames appends the other names to the license details.
func (license License) appendOtherNames(builder *strings.Builder, bold *color.Color) {
	if license.OtherNames != nil && len(*license.OtherNames) > 0 {
		builder.WriteString(bold.Sprintln("\nCommon names"))
		for _, other := range *license.OtherNames {
			builder.WriteString(fmt.Sprintf("  * %s\n", other.Name))
		}
	}
}

// appendIdentifiers appends the identifiers to the license details.
func (license License) appendIdentifiers(builder *strings.Builder, bold *color.Color) {
	if len(license.Identifiers) > 0 {
		builder.WriteString(bold.Sprintln("\nLicense Standards"))
		for _, identifier := range license.Identifiers {
			builder.WriteString(fmt.Sprintf("  * %-10s %s\t\n", identifier.Scheme, identifier.Identifier))
		}
	}
}

// appendLinks appends the links to the license details.
func (license License) appendLinks(builder *strings.Builder, bold *color.Color, italic *color.Color) {
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
}

// FindById enable searching for a specific license
func (l Licenses) FindById(id string) *License {
	for _, license := range l {
		if strings.Contains(strings.ToLower(license.Id), strings.ToLower(id)) {
			return &license
		}
	}
	return nil
}

// FindByKeyword enable querying by keyword (e.g. "popular")
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
	for _, license := range l {
		if l.matchesLicense(license, term) {
			licenses = append(licenses, license)
		}
	}
	return &licenses
}

// matchesLicense checks if a license matches the search term
func (l Licenses) matchesLicense(license License, term string) bool {
	if strings.Contains(strings.ToLower(license.Id), strings.ToLower(term)) {
		return true
	}
	if l.matchesIdentifiers(license.Identifiers, term) {
		return true
	}
	if fuzzy.MatchNormalizedFold(term, license.Name) {
		return true
	}
	if l.matchesOtherNames(license.OtherNames, term) {
		return true
	}
	return false
}

// matchesIdentifiers checks if a license matches the search term
func (l Licenses) matchesIdentifiers(identifiers []Identifier, term string) bool {
	for _, identifier := range identifiers {
		if identifier.Scheme == "SPDX" && term == identifier.Identifier {
			return true
		}
		if fuzzy.MatchNormalizedFold(term, identifier.Identifier) {
			return true
		}
	}
	return false
}

// matchesOtherNames checks if a license matches the search term
func (l Licenses) matchesOtherNames(otherNames *[]OtherName, term string) bool {
	if otherNames == nil {
		return false
	}
	for _, otherName := range *otherNames {
		if fuzzy.MatchNormalizedFold(term, otherName.Name) {
			return true
		}
	}
	return false
}
