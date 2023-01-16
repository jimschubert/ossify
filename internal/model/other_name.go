package model

type OtherName struct {
	// pointer works for missing parameters, verify this works for null values.
	Note *string `json:"note"`
	Name string  `json:"name"`
}
