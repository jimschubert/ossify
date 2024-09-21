package model

// Text represents a text resource with a title, URL, and media type.
type Text struct {
	// Title is the title of the text resource.
	Title string `json:"title"`

	// Url is the URL where the text resource can be accessed.
	Url string `json:"url"`

	// MediaType is the type of media of the text resource (e.g., text/html, text/plain).
	MediaType string `json:"media_type"`
}
