package model

import (
	"fmt"
	"strings"
)

type License struct {
	SupersededBy *string `json:"superseded_by"`
	Identifiers []Identifier `json:"identifiers"`
	Text []Text `json:"text"`
	OtherNames *[]OtherName `json:"other_names"`
	Id string `json:"id"`
	Links []Link `json:"links"`
	Name string `json:"name"`
	Keywords []string `json:"keywords"`
}

func (license License) Print() error {
	_, err := fmt.Printf("%-20s(%s)\n", license.Id, license.Name)
	return err
}

// This type alias allows us to create a sort of "dao" atop the set of data.
// see Go in Action Chapter 2 for somewhat similar approach.
type Licenses []License

// enable searching for a specific license
func (l Licenses) FindById(id string) *License {
	for _, license := range l {
		if license.Id == id {
			return &license
		}
	}
	return nil
}

// enable querying by keyword (e.g. "popular")
func (l Licenses) FindByKeyword(keyword string) *Licenses {
	var licenses Licenses
	for _, license := range l {
		if license.Id == keyword {
			licenses = append(licenses, license)
		}
	}
	return &licenses
}

// enable a loose free-form textual search
func (l Licenses) Search(term string) *Licenses {
	var licenses Licenses
	for _, license := range l {
		var added = false
		if strings.Contains(license.Id, term) {
			licenses = append(licenses, license)
			added = true
		}

		if !added && license.Identifiers != nil {
			for _, identifier := range license.Identifiers {
				if added {
					break
				}
				if strings.Contains(identifier.Identifier, term) {
					licenses = append(licenses, license)
					added = true
				}
			}
		}

		if !added && license.Text != nil {
			for _, text := range license.Text {
				if added {
					break
				}
				if strings.Contains(text.Title, term) {
					licenses = append(licenses, license)
					added = true
				}
			}
		}

		if !added {
			if strings.Contains(license.Name, term) {
				licenses = append(licenses, license)
				added = true
			}
		}

		if !added && license.OtherNames != nil {
			for _, otherNames := range *license.OtherNames {
				if added {
					break
				}
				if otherNames.Note != nil && strings.Contains(*otherNames.Note, term) {
					licenses = append(licenses, license)
					added = true
				}

				if !added && strings.Contains(otherNames.Name, term) {
					licenses = append(licenses, license)
					added = true
				}
			}
		}
	}
	return &licenses
}