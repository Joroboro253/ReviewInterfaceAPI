package models

type ReviewUpdateRequest struct {
	Data struct {
		Type       string       `json:"type"`
		Attributes ReviewUpdate `json:"attributes"`
	} `json:"data"`
}
