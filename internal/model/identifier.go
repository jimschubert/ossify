package model

// Identifier represents a unique identifier within a specific scheme.
type Identifier struct {
	// Scheme is the name of the scheme to which the identifier belongs.
	Scheme string `json:"scheme"`

	// Identifier is the unique identifier within the specified scheme.
	Identifier string `json:"identifier"`
}
