package models

type ReviewUpdateData struct {
	Type       string       `json:"type"`
	ID         string       `json:"id"`
	Attributes ReviewUpdate `json:"attributes"`
}

type ReviewUpdateRequest struct {
	Data ReviewUpdateData `json:"data"`
}
