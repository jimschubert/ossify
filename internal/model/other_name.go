package model

// OtherName represents an alternative name for a license.
// It includes a note that can be optionally provided.
type OtherName struct {
	// Note is an optional field that provides additional information about the name.
	// It is a pointer to a string to handle missing or null values.
	Note *string `json:"note"`

	// Name is the alternative name for the license.
	Name string `json:"name"`
}
