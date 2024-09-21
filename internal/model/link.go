package model

// Link represents a hyperlink with an optional note.
// It includes a URL and an optional note that provides additional information about the link.
type Link struct {
	// Note is an optional field that provides additional information about the link.
	// It is a pointer to a string to handle missing or null values.
	Note *string `json:"note"`

	// Url is the URL of the hyperlink.
	Url string `json:"url"`
}
