package models

type ReviewData struct {
	Type       string `json:"type,omitempty"`
	Attributes Review `json:"attributes,omitempty"`
}
